# Getting Started with Keycloak (Local Development)

This guide explains how to run Keycloak locally using this project's Docker setup. Keycloak is auto-provisioned on startup with a realm, client, and users so you don't have to click around.

## Prerequisites
- Docker and Docker Compose installed
- This repository cloned locally

## Files involved
- docker/docker-compose.yml – defines the keycloak service
- docker/keycloak.defaults.env – non-sensitive defaults (checked in)
- docker/keycloak.user.env – your local admin credentials (checked in with placeholders; set values locally)
- docker/keycloak.user.env.example – reference example

## 1) Set bootstrap admin credentials
Edit `docker/keycloak.user.env` and set both variables (used by Keycloak to start the server). The file is checked in with BLANK values by default, so you must provide your own locally before starting:

```
KEYCLOAK_ADMIN=your-username
KEYCLOAK_ADMIN_PASSWORD=your-strong-password
```

Keycloak will not start unless both variables are defined. These are the bootstrap Admin Console credentials.

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

On startup, a realm called `traveler-dev` is automatically imported with a preconfigured client and users (see below). Note: This dev setup does not persist Keycloak data to a volume. Removing the container will remove realms, clients, and users; they will be re-imported next start.

## 3) What's provisioned automatically
- Realm: `traveler-dev`
- Client: `traveler-app` (OIDC public)
  - Redirect URIs: `http://localhost:8080/*`
  - Web origins: `http://localhost:8080`
- Roles:
  - Realm role: `api-user`
- Users:
  - Admin user: `kc_admin2` with password `Admin#1!` (has `realm-management -> realm-admin` role; can administer the realm)
  - API user: `api-user` with password `ApiUser#1!` (granted realm role `api-user`)

You can log in to the Admin Console either with the bootstrap admin (from `docker/keycloak.user.env`) or with the pre-provisioned `kc_admin2` above.

## 4) Get a token (quick CLI test)
Replace placeholders and run:

```
curl -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=traveler-app" \
  -d "username=api-user" \
  -d "password=ApiUser#1!" \
  http://localhost:8081/realms/traveler-dev/protocol/openid-connect/token
```

Well-known config for your realm:

```
http://localhost:8081/realms/traveler-dev/.well-known/openid-configuration
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
