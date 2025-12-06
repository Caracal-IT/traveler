# Getting Started with Keycloak (Local Development)

This guide explains how to run Keycloak locally using this project's Docker setup and perform minimal admin steps for testing.

## Prerequisites
- Docker and Docker Compose installed
- This repository cloned locally

## Files involved
- docker/docker-compose.yml – defines the keycloak service
- docker/keycloak.defaults.env – non-sensitive defaults (checked in)
- docker/keycloak.user.env – your local admin credentials (checked in with placeholders; set values locally)
- docker/keycloak.user.env.example – reference example

## 1) Set admin credentials
Edit docker/keycloak.user.env and set both variables:

```
KEYCLOAK_ADMIN=your-username
KEYCLOAK_ADMIN_PASSWORD=your-strong-password
```

Keycloak will not start unless both variables are defined.

## 2) Start services
From the project root:

```
cd docker
docker compose up -d --build
```

Services:
- Traveler app: http://localhost:8080
- Keycloak: http://localhost:8081

Open the Keycloak Admin Console at http://localhost:8081 and log in with the admin credentials you set.

Note: This dev setup does not persist Keycloak data to a volume. Removing the container will remove realms, clients, and users.

## 3) Create a realm
1. In the Admin Console, open the realm dropdown (top-left)
2. Click Create realm
3. Name it (for example, traveler-dev) and save

## 4) Create a client (OIDC)
1. Inside your realm, go to Clients -> Create client
2. Client ID: traveler-app (or your choice)
3. Protocol: OpenID Connect -> Next
4. Configure:
   - Valid redirect URIs: http://localhost:8080/* (for local testing)
   - Web origins: http://localhost:8080 (or * for broad local testing only)
5. Save

## 5) Create a test user
1. Users -> Add user -> set username (for example, alice) -> Save
2. Open the user -> Credentials -> Set password
3. Provide a password and disable Temporary if you do not want to force a reset

## 6) Get a token (quick CLI test)
Replace placeholders and run:

```
curl -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=<client_id>" \
  -d "username=<username>" \
  -d "password=<password>" \
  http://localhost:8081/realms/<realm>/protocol/openid-connect/token
```

Well-known config for your realm:

```
http://localhost:8081/realms/<realm>/.well-known/openid-configuration
```

## Troubleshooting
- Ensure both KEYCLOAK_ADMIN and KEYCLOAK_ADMIN_PASSWORD are set
- Port conflict: change the host port mapping in docker/docker-compose.yml
- View logs: from the docker directory, run: docker compose logs -f keycloak
- Reset to a clean state: from the docker directory, run: docker compose down -v (deletes data)

## Stop services
```
cd docker
docker compose down
```
