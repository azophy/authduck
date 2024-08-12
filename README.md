<p align="center">
  <img width="200" height="200" src="./resources/assets/logo.svg?raw=true" alt="Authduck Logo" />
</p>

# AuthDuck
OIDC Server Playground for development & testing

## Problem Statement

Sometime when developing a web app, its easier to integrate OIDC early in the process, as it allows to skip many important auth-related feature (login, registration, password reset, etc). however setting up IdP for development purposes often takes a lot of effort. this project aims to provide a simple IdP server for this purposes, with minimal features

registering on a cloud OIDC provider is not always straightforward (google, twitter, etc). distributing the credentials is also another minefield. so a mock IdP server seems like the best option

## How to use

Curently we only provide a generic OIDC server module. However in the future we plan to add mock server for popular services such as Google, Github, Keycloak, etc.

To start playing with this tools as a Generic OIDC Server, you could use any OIDC client you have (for example: <a href="https://openidconnect.net" target="_blank">openidconnect.net</a> , an OIDC client playground by <a href="https://auth0.com" target="_blank">Auth0</a>), and then use the <a href="https://authduck.fly.dev/case/generic/.well-known/openid-configuration">/.well-known/openid-configuration</a> path to populate all the required fields. You could also view request history to known client by visiting the <a target="_blank" href="./manage/history">request history page</a>.

There are a [tutorial in our wiki page](https://github.com/azophy/authduck/wiki/How-to-play-with-Authduck-&-Openidconnect.net) if you wish for a more detailed guidance.

## How To Install

For many cases, using our public instance at <a target="_blank" href="https://authduck.fly.dev">authduck.fly.dev</a> should be enough. However if you wish to install it on your own server, we provide several methods:

### 1. Download available executables from releases page

> Due to complications on cross-compiling go with CGO enabled, we currently only provide linux-amd64 distribution

### 2. Use available docker image

```
docker run -it ghcr.io/azophy/authduck:latest
```

### 3. Compile from source

commonly for development. just prepare go version 1.22+ and run `go build`.

## Features
- [x] well known endpoint
- [x] jwks endpoint
- [x] customize call back payload
- [ ] define accepted clients & users
- [x] deployable as single binary
- [x] dockerized in dockerhub
- [ ] rate limiter to protect from abuse
- [ ] scheduler to clear stale data
- [ ] mock for popular OIDC providers: google, facebook, github, keycloak, ory, etc

## Tech stack
- go 1.22
- in memory sqlite (configurable)
- tacit css framework
- HTMX
- JWX package for JWT, JWK, and JWS operations


