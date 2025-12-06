package auth

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"traveler/pkg/config"
	"traveler/pkg/log"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwksOnce sync.Once
	jwksMap  = make(map[string]*keyfunc.JWKS)
	jwksErr  error
	mu       sync.RWMutex
)

// getJWKS returns a cached JWKS for the given JWKS URL.
func getJWKS(jwksURL string) (*keyfunc.JWKS, error) {
	mu.RLock()
	if jwks, ok := jwksMap[jwksURL]; ok && jwks != nil {
		mu.RUnlock()
		return jwks, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()
	if jwks, ok := jwksMap[jwksURL]; ok && jwks != nil {
		return jwks, nil
	}

	options := keyfunc.Options{}
	// Enable background refresh with sane intervals
	options.RefreshErrorHandler = func(err error) {
		log.Warn("JWKS refresh error", "error", err)
	}
	options.RefreshInterval = time.Hour
	options.RefreshTimeout = 5 * time.Second
	options.Client = &http.Client{Timeout: 5 * time.Second}

	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		return nil, err
	}
	jwksMap[jwksURL] = jwks
	return jwks, nil
}

// JWTMiddleware validates Bearer tokens issued by Keycloak and enforces issuer/audience.
func JWTMiddleware(cfg *config.Config) fiber.Handler {
	issuer := cfg.Auth.Issuer
	audience := cfg.Auth.Audience
	// Compute JWKS URL — allow override via config to support containerized envs where issuer host differs
	jwksURL := strings.TrimRight(issuer, "/") + "/protocol/openid-connect/certs"
	if cfg.Auth.JWKSURL != "" {
		jwksURL = cfg.Auth.JWKSURL
	}

	return func(c *fiber.Ctx) error {
		authz := c.Get("Authorization")
		if authz == "" {
			return fiber.ErrUnauthorized
		}
		parts := strings.SplitN(authz, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return fiber.ErrUnauthorized
		}
		tokenString := parts[1]

		jwks, err := getJWKS(jwksURL)
		if err != nil {
			log.Error("failed to get JWKS", "error", err)
			return fiber.ErrUnauthorized
		}

		// First, validate signature, issuer, and algorithm. We'll handle audience
		// validation manually to be compatible with Keycloak where `aud` may be
		// "account" and the client id appears in `azp` (authorized party).
		parsed, err := jwt.Parse(
			tokenString,
			jwks.Keyfunc,
			// Validate signature algorithm; we'll validate issuer manually for better dev flexibility
			jwt.WithValidMethods([]string{"RS256", "RS384", "RS512"}),
			// Allow small clock skew to avoid 401 due to exp/nbf/iat drift in local/dev environments
			jwt.WithLeeway(60*time.Second),
		)
		if err != nil || !parsed.Valid {
			if err == nil {
				err = errors.New("invalid token")
			}
			log.Warn("token validation failed", "error", err)
			return fiber.ErrUnauthorized
		}

		// Issuer validation (done manually to tolerate http/https & trailing slash differences in dev)
		if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
			issClaim, _ := claims["iss"].(string)
			if !issuerAllowed(issuer, issClaim) {
				log.Warn("token issuer mismatch", "expected_issuer", issuer, "token_iss", issClaim)
				return fiber.ErrUnauthorized
			}

			// Audience/Client validation compatible with Keycloak
			// 1) Try standard aud claim (string or array)
			audOK := false
			if rawAud, exists := claims["aud"]; exists {
				switch v := rawAud.(type) {
				case string:
					if strings.EqualFold(v, audience) {
						audOK = true
					}
				case []interface{}:
					for _, item := range v {
						if s, ok := item.(string); ok && strings.EqualFold(s, audience) {
							audOK = true
							break
						}
					}
				}
			}

			// 2) Fallback to Keycloak's `azp` (authorized party == client_id)
			if !audOK {
				if azp, ok := claims["azp"].(string); ok && strings.EqualFold(azp, audience) {
					audOK = true
				}
			}

			// 3) Fallback to Keycloak's `resource_access` map which lists client roles
			//    Accept token if it contains an entry for our client id regardless of roles
			if !audOK {
				if ra, ok := claims["resource_access"].(map[string]interface{}); ok {
					if _, exists := ra[audience]; exists {
						audOK = true
					}
				}
			}

			if !audOK {
				log.Warn("token audience/azp mismatch", "expected_audience", audience, "claims_aud", claims["aud"], "claims_azp", claims["azp"])
				return fiber.ErrUnauthorized
			}

			// Store token claims in context for handlers to use
			c.Locals("claims", claims)
		}
		return c.Next()
	}
}

// issuerAllowed compares expected issuer with token issuer allowing small
// variations common in local development (http vs https and trailing slash).
func issuerAllowed(expected, actual string) bool {
	if expected == "" || actual == "" {
		return false
	}
	norm := func(s string) string {
		return strings.TrimRight(s, "/")
	}
	e := norm(expected)
	a := norm(actual)
	if strings.EqualFold(e, a) {
		return true
	}
	// Allow http/https swap for local dev
	swapScheme := func(u string) string {
		if strings.HasPrefix(u, "http://") {
			return "https://" + strings.TrimPrefix(u, "http://")
		}
		if strings.HasPrefix(u, "https://") {
			return "http://" + strings.TrimPrefix(u, "https://")
		}
		return u
	}
	if strings.EqualFold(swapScheme(e), a) {
		return true
	}
	return false
}
