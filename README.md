# AuthDuck
OIDC Identity Provider for development & testing

## Problem Statement 

Sometime when developing a web app, its easier to integrate OIDC early in the process, as it allows to skip many important auth-related feature (login, registration, password reset, etc). however setting up IdP for development purposes often takes a lot of effort. this project aims to provide a simple IdP server for this purposes, with minimal features

registering on a cloud OIDC provider is not always straightforward (google, twitter, etc). distributing the credentials is also another minefield. so local IdP server seems like the best option

## Features
- well known endpoint
- jwks endpoint
- customize call back payload
- define accepted clients & users
- dockerized in dockerhub
