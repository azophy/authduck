![Authduck Logo](resources/assets/logo.png?raw=true "Authduck Logo")

# AuthDuck
OIDC Identity Provider for development & testing

## Problem Statement

Sometime when developing a web app, its easier to integrate OIDC early in the process, as it allows to skip many important auth-related feature (login, registration, password reset, etc). however setting up IdP for development purposes often takes a lot of effort. this project aims to provide a simple IdP server for this purposes, with minimal features

registering on a cloud OIDC provider is not always straightforward (google, twitter, etc). distributing the credentials is also another minefield. so local IdP server seems like the best option

## Features
- [x] well known endpoint
- [x] jwks endpoint
- [x] customize call back payload
- [ ] define accepted clients & users
- [ ] deployable as single binary
- [ ] dockerized in dockerhub
- [ ] rate limiter to protect from abuse
- [ ] scheduler to clear stale data
- [ ] mock for popular OIDC providers: google, facebook, github, keycloak, ory, etc

## Tech stack
- go 1.22
- in memory sqlite (configurable)
- tacit css framework
- HTMX
- JWX package for JWT, JWK, and JWS operations


