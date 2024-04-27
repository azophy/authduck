package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})


	e.GET("/.well-known/openid-configuration", func(c echo.Context) error {
		return c.String(http.StatusOK, `{
      "issuer": "https://example.com/realms/ssojabar",
      "authorization_endpoint": "https://example.com/protocol/openid-connect/auth",
      "token_endpoint": "https://example.com/protocol/openid-connect/token",
      "introspection_endpoint": "https://example.com/protocol/openid-connect/token/introspect",
      "userinfo_endpoint": "https://example.com/protocol/openid-connect/userinfo",
      "end_session_endpoint": "https://example.com/protocol/openid-connect/logout",
      "frontchannel_logout_session_supported": true,
      "frontchannel_logout_supported": true,
      "jwks_uri": "https://example.com/protocol/openid-connect/certs",
      "check_session_iframe": "https://example.com/protocol/openid-connect/login-status-iframe.html",
      "grant_types_supported": [
        "authorization_code",
        "implicit",
        "refresh_token",
        "password",
        "client_credentials",
        "urn:ietf:params:oauth:grant-type:device_code",
        "urn:openid:params:grant-type:ciba"
      ],
      "response_types_supported": [
        "code",
        "none",
        "id_token",
        "token",
        "id_token token",
        "code id_token",
        "code token",
        "code id_token token"
      ],
      "subject_types_supported": [
        "public",
        "pairwise"
      ],
      "id_token_signing_alg_values_supported": [
        "PS384",
        "ES384",
        "RS384",
        "HS256",
        "HS512",
        "ES256",
        "RS256",
        "HS384",
        "ES512",
        "PS256",
        "PS512",
        "RS512"
      ],
      "id_token_encryption_alg_values_supported": [
        "RSA-OAEP",
        "RSA-OAEP-256",
        "RSA1_5"
      ],
      "id_token_encryption_enc_values_supported": [
        "A256GCM",
        "A192GCM",
        "A128GCM",
        "A128CBC-HS256",
        "A192CBC-HS384",
        "A256CBC-HS512"
      ],
      "userinfo_signing_alg_values_supported": [
        "PS384",
        "ES384",
        "RS384",
        "HS256",
        "HS512",
        "ES256",
        "RS256",
        "HS384",
        "ES512",
        "PS256",
        "PS512",
        "RS512",
        "none"
      ],
      "request_object_signing_alg_values_supported": [
        "PS384",
        "ES384",
        "RS384",
        "HS256",
        "HS512",
        "ES256",
        "RS256",
        "HS384",
        "ES512",
        "PS256",
        "PS512",
        "RS512",
        "none"
      ],
      "request_object_encryption_alg_values_supported": [
        "RSA-OAEP",
        "RSA-OAEP-256",
        "RSA1_5"
      ],
      "request_object_encryption_enc_values_supported": [
        "A256GCM",
        "A192GCM",
        "A128GCM",
        "A128CBC-HS256",
        "A192CBC-HS384",
        "A256CBC-HS512"
      ],
      "response_modes_supported": [
        "query",
        "fragment",
        "form_post",
        "query.jwt",
        "fragment.jwt",
        "form_post.jwt",
        "jwt"
      ],
      "registration_endpoint": "https://example.com/clients-registrations/openid-connect",
      "token_endpoint_auth_methods_supported": [
        "private_key_jwt",
        "client_secret_basic",
        "client_secret_post",
        "tls_client_auth",
        "client_secret_jwt"
      ],
      "token_endpoint_auth_signing_alg_values_supported": [
        "PS384",
        "ES384",
        "RS384",
        "HS256",
        "HS512",
        "ES256",
        "RS256",
        "HS384",
        "ES512",
        "PS256",
        "PS512",
        "RS512"
      ],
      "introspection_endpoint_auth_methods_supported": [
        "private_key_jwt",
        "client_secret_basic",
        "client_secret_post",
        "tls_client_auth",
        "client_secret_jwt"
      ],
      "introspection_endpoint_auth_signing_alg_values_supported": [
        "PS384",
        "ES384",
        "RS384",
        "HS256",
        "HS512",
        "ES256",
        "RS256",
        "HS384",
        "ES512",
        "PS256",
        "PS512",
        "RS512"
      ],
      "authorization_signing_alg_values_supported": [
        "PS384",
        "ES384",
        "RS384",
        "HS256",
        "HS512",
        "ES256",
        "RS256",
        "HS384",
        "ES512",
        "PS256",
        "PS512",
        "RS512"
      ],
      "authorization_encryption_alg_values_supported": [
        "RSA-OAEP",
        "RSA-OAEP-256",
        "RSA1_5"
      ],
      "authorization_encryption_enc_values_supported": [
        "A256GCM",
        "A192GCM",
        "A128GCM",
        "A128CBC-HS256",
        "A192CBC-HS384",
        "A256CBC-HS512"
      ],
      "claims_supported": [
        "aud",
        "sub",
        "iss",
        "auth_time",
        "name",
        "given_name",
        "family_name",
        "preferred_username",
        "email",
        "acr"
      ],
      "claim_types_supported": [
        "normal"
      ],
      "claims_parameter_supported": true,
      "scopes_supported": [
        "openid",
        "roles",
        "address",
        "profile",
        "acr",
        "email",
        "phone",
        "offline_access",
        "web-origins",
        "microprofile-jwt"
      ],
      "request_parameter_supported": true,
      "request_uri_parameter_supported": true,
      "require_request_uri_registration": true,
      "code_challenge_methods_supported": [
        "plain",
        "S256"
      ],
      "tls_client_certificate_bound_access_tokens": true,
      "revocation_endpoint": "https://example.com/protocol/openid-connect/revoke",
      "revocation_endpoint_auth_methods_supported": [
        "private_key_jwt",
        "client_secret_basic",
        "client_secret_post",
        "tls_client_auth",
        "client_secret_jwt"
      ],
      "revocation_endpoint_auth_signing_alg_values_supported": [
        "PS384",
        "ES384",
        "RS384",
        "HS256",
        "HS512",
        "ES256",
        "RS256",
        "HS384",
        "ES512",
        "PS256",
        "PS512",
        "RS512"
      ],
      "backchannel_logout_supported": true,
      "backchannel_logout_session_supported": true,
      "device_authorization_endpoint": "https://example.com/protocol/openid-connect/auth/device",
      "backchannel_token_delivery_modes_supported": [
        "poll",
        "ping"
      ],
      "backchannel_authentication_endpoint": "https://example.com/protocol/openid-connect/ext/ciba/auth",
      "backchannel_authentication_request_signing_alg_values_supported": [
        "PS384",
        "ES384",
        "RS384",
        "ES256",
        "RS256",
        "ES512",
        "PS256",
        "PS512",
        "RS512"
      ],
      "require_pushed_authorization_requests": false,
      "pushed_authorization_request_endpoint": "https://example.com/protocol/openid-connect/ext/par/request",
      "mtls_endpoint_aliases": {
        "token_endpoint": "https://example.com/protocol/openid-connect/token",
        "revocation_endpoint": "https://example.com/protocol/openid-connect/revoke",
        "introspection_endpoint": "https://example.com/protocol/openid-connect/token/introspect",
        "device_authorization_endpoint": "https://example.com/protocol/openid-connect/auth/device",
        "registration_endpoint": "https://example.com/clients-registrations/openid-connect",
        "userinfo_endpoint": "https://example.com/protocol/openid-connect/userinfo",
        "pushed_authorization_request_endpoint": "https://example.com/protocol/openid-connect/ext/par/request",
        "backchannel_authentication_endpoint": "https://example.com/protocol/openid-connect/ext/ciba/auth"
      }
    }`)
	})

	e.Logger.Fatal(e.Start(":3000"))
}