# `radau` - Radius authentication backend

Radius Authentication REST backend microservice used to manage users and credentials for authentication in a WPA Enterprise setup. The credentials are ready to be passed to a radius server, e.g. the FreeRadius `rest` module.

[📝 API docs](https://factorycampus.github.io/radau/)

**Important:** Passwords can not be hashed in the database to allow comparison to various authentication methods supported in different clients.

The implementation therefore adheres to these principles:
- Passwords are always generated
- The user's password is only returned on creation
- Passwords can be reset or revoked

### Additional features
The API allows for additional attributes to be associated to the user which will be passed to the radius server on authentication

## Usage

### Requirements
- A PostgreSQL database

### Docker

```bash
docker run -d --name wifilogin_db --network wifi_db \
    -e POSTGRES_USER=wifi -e POSTGRES_PASSWORD=wifi -e POSTGRES_DB=wifi \
    postgres
docker run -d --name wifilogin \
    -e API_KEY_MANAGEMENT= -e API_KEY_RADIUS= \
    -e DB_HOST=wifilogin_db -e DB_USER=wifi -e DB_PASSWORD=wifi -e DB_DATABASE=wifi \
    -p 8080:8080 --network wifi_db
    factorycampus/wifi-login-backend
```

### Debian

A debian package is available on the GitHub release page.

### Available Environment Variables

- `API_KEY_MANAGEMENT` - API Key used for user management (`/user` and `/token`)
- `API_KEY_RADIUS` - API Key used by radius server (`/radius`)
- `DB_HOST` - Host of the database server
- `DB_USER` - User for database access
- `DB_PASSWORD` - Password for database access
- `DB_DATABASE` - Database to use
- `PORT` - Port to serve the API on
- `TOKEN_LENGTH` - Length of the generated token, defaults to 32

### In Production

Consider these tips:

- Because credentials will be in plaintext, put this service behind a SSL-enabled reverse-proxy and access via HTTPS from Radius

## Development

Use the `docker-compose.yml` in the repo root for development with livereload
