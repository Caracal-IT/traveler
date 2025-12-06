Specials API
===========

Endpoint to retrieve a list of special travel packages. This endpoint is protected by Keycloak and requires a valid Bearer token.

- Method: GET
- URL: /api/offerings/specials
- Auth: Bearer token (JWT) issued by Keycloak realm `traveler-dev` for audience `traveler-app`

Authentication
--------------
- Issuer: http://localhost:8081/realms/traveler-dev
- Audience (client_id): traveler-app
- Example user (local dev realm): api-user / ApiUser#1!

Response
--------
- 200 OK

```
{
  "items": [
    {"id": "sp-1001", "name": "Winter Escape", "price": 799.0, "currency": "USD"},
    {"id": "sp-1002", "name": "City Break Deluxe", "price": 499.0, "currency": "USD"}
  ]
}
```

- 401 Unauthorized – missing/invalid token
- 403 Forbidden – token valid but not permitted

Quick test (cURL)
-----------------
1) Get a token

```
ACCESS_TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=traveler-app" \
  -d "username=api-user" \
  -d "password=ApiUser#1!" \
  http://localhost:8081/realms/traveler-dev/protocol/openid-connect/token | jq -r .access_token)
```

2) Call the endpoint

```
curl -H "Authorization: Bearer $ACCESS_TOKEN" \
  http://localhost:8080/api/offerings/specials | jq
```

OpenAPI
-------
See api/openapi.yaml under path /api/offerings/specials with bearerAuth security.
