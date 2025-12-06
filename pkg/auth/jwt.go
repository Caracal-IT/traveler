package auth

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"traveler/pkg/config"
	"traveler/pkg/log"
)

var (
	jwksOnce sync.Once
	jwksMap  = make(map[string]*keyfunc.JWKS)
	jwksErr  error
	mu       sync.RWMutex
)

// getJWKS returns a cached JWKS for the given issuer.
func getJWKS(issuer string) (*keyfunc.JWKS, error) {
	jwksURL := strings.TrimRight(issuer, "/") + "/protocol/openid-connect/certs"
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

		jwks, err := getJWKS(issuer)
		if err != nil {
			log.Error("failed to get JWKS", "error", err)
			return fiber.ErrUnauthorized
		}

		parsed, err := jwt.Parse(tokenString, jwks.Keyfunc,
			jwt.WithAudience(audience),
			jwt.WithIssuer(issuer),
			jwt.WithValidMethods([]string{"RS256", "RS384", "RS512"}),
		)
		if err != nil || !parsed.Valid {
			if err == nil {
				err = errors.New("invalid token")
			}
			log.Warn("token validation failed", "error", err)
			return fiber.ErrUnauthorized
		}

		// Store token claims in context for handlers to use
		if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
			c.Locals("claims", claims)
		}
		return c.Next()
	}
}
