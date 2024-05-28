package main

import (
  "log"
  "time"
  "net/http"
  "math/rand"

	"github.com/labstack/echo/v4"
  "github.com/lestrrat-go/jwx/v2/jwa"
  "github.com/lestrrat-go/jwx/v2/jwt"
  "github.com/lestrrat-go/jwx/v2/jwt/openid"
)

const (
  routeParent = "/case/general"
)

func RegisterGeneralOAuthModule(app *echo.Echo) {
  e := app.Group(routeParent)

	e.GET("/.well-known/openid-configuration", openidconfigHandler)
	//e.GET("/auth/callback", CallbackHandler)
  e.GET("/auth/callback", ServeResourceTemplate("resources/views/generic_callback.html", nil))
	e.POST("/auth/token", tokenHandler)
}

func openidconfigHandler(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
      "issuer": BaseUrl,
      "authorization_endpoint": BaseUrl + routeParent + "/auth/callback",
      "token_endpoint": BaseUrl + routeParent + "/auth/token",
      "introspection_endpoint": BaseUrl + routeParent + "/auth/introspect",
      "userinfo_endpoint": BaseUrl + routeParent + "/auth/userinfo",
      "end_session_endpoint": BaseUrl + routeParent + "/auth/logout",
      //"frontchannel_logout_session_supported": true,
      //"frontchannel_logout_supported": true,
      "jwks_uri": BaseUrl + "/.well-known/certs",
      //"check_session_iframe": BaseUrl + "/protocol/openid-connect/login-status-iframe.html",
      "grant_types_supported": []string{
        "authorization_code",
        //"implicit",
        //"refresh_token",
        //"password",
        //"client_credentials",
        //"urn:ietf:params:oauth:grant-type:device_code",
        //"urn:openid:params:grant-type:ciba",
      },
      //"response_types_supported": []string{
        //"code",
        //"none",
        //"id_token",
        //"token",
        //"id_token token",
        //"code id_token",
        //"code token",
        //"code id_token token",
      //},
      //"subject_types_supported": []string{
        //"public",
        //"pairwise",
      //},
      //"id_token_signing_alg_values_supported": []string{
        //"PS384",
        //"ES384",
        //"RS384",
        //"HS256",
        //"HS512",
        //"ES256",
        //"RS256",
        //"HS384",
        //"ES512",
        //"PS256",
        //"PS512",
        //"RS512",
      //},
      //"id_token_encryption_alg_values_supported": []string{
        //"RSA-OAEP",
        //"RSA-OAEP-256",
        //"RSA1_5",
      //},
      //"id_token_encryption_enc_values_supported": []string{
        //"A256GCM",
        //"A192GCM",
        //"A128GCM",
        //"A128CBC-HS256",
        //"A192CBC-HS384",
        //"A256CBC-HS512",
      //},
      //"userinfo_signing_alg_values_supported": []string{
        //"PS384",
        //"ES384",
        //"RS384",
        //"HS256",
        //"HS512",
        //"ES256",
        //"RS256",
        //"HS384",
        //"ES512",
        //"PS256",
        //"PS512",
        //"RS512",
        //"none",
      //},
      //"request_object_signing_alg_values_supported": []string{
        //"PS384",
        //"ES384",
        //"RS384",
        //"HS256",
        //"HS512",
        //"ES256",
        //"RS256",
        //"HS384",
        //"ES512",
        //"PS256",
        //"PS512",
        //"RS512",
        //"none",
      //},
      //"request_object_encryption_alg_values_supported": []string{
        //"RSA-OAEP",
        //"RSA-OAEP-256",
        //"RSA1_5",
      //},
      //"request_object_encryption_enc_values_supported": []string{
        //"A256GCM",
        //"A192GCM",
        //"A128GCM",
        //"A128CBC-HS256",
        //"A192CBC-HS384",
        //"A256CBC-HS512",
      //},
      //"response_modes_supported": []string{
        //"query",
        //"fragment",
        //"form_post",
        //"query.jwt",
        //"fragment.jwt",
        //"form_post.jwt",
        //"jwt",
      //},
      //"registration_endpoint": "https://example.com/clients-registrations/openid-connect",
      //"token_endpoint_auth_methods_supported": []string{
        //"private_key_jwt",
        //"client_secret_basic",
        //"client_secret_post",
        //"tls_client_auth",
        //"client_secret_jwt",
      //},
      //"token_endpoint_auth_signing_alg_values_supported": []string{
        //"PS384",
        //"ES384",
        //"RS384",
        //"HS256",
        //"HS512",
        //"ES256",
        //"RS256",
        //"HS384",
        //"ES512",
        //"PS256",
        //"PS512",
        //"RS512",
      //},
      //"introspection_endpoint_auth_methods_supported": []string{
        //"private_key_jwt",
        //"client_secret_basic",
        //"client_secret_post",
        //"tls_client_auth",
        //"client_secret_jwt",
      //},
      //"introspection_endpoint_auth_signing_alg_values_supported": []string{
        //"PS384",
        //"ES384",
        //"RS384",
        //"HS256",
        //"HS512",
        //"ES256",
        //"RS256",
        //"HS384",
        //"ES512",
        //"PS256",
        //"PS512",
        //"RS512",
      //},
      //"authorization_signing_alg_values_supported": []string{
        //"PS384",
        //"ES384",
        //"RS384",
        //"HS256",
        //"HS512",
        //"ES256",
        //"RS256",
        //"HS384",
        //"ES512",
        //"PS256",
        //"PS512",
        //"RS512",
      //},
      //"authorization_encryption_alg_values_supported": []string{
        //"RSA-OAEP",
        //"RSA-OAEP-256",
        //"RSA1_5",
      //},
      //"authorization_encryption_enc_values_supported": []string{
        //"A256GCM",
        //"A192GCM",
        //"A128GCM",
        //"A128CBC-HS256",
        //"A192CBC-HS384",
        //"A256CBC-HS512",
      //},
      //"claims_supported": []string{
        //"aud",
        //"sub",
        //"iss",
        //"auth_time",
        //"name",
        //"given_name",
        //"family_name",
        //"preferred_username",
        //"email",
        //"acr",
      //},
      //"claim_types_supported": []string{
        //"normal",
      //},
      //"claims_parameter_supported": true,
      //"scopes_supported": []string{
        //"openid",
        //"roles",
        //"address",
        //"profile",
        //"acr",
        //"email",
        //"phone",
        //"offline_access",
        //"web-origins",
        //"microprofile-jwt",
      //},
      //"request_parameter_supported": true,
      //"request_uri_parameter_supported": true,
      //"require_request_uri_registration": true,
      //"code_challenge_methods_supported": []string{
        //"plain",
        //"S256",
      //},
      //"tls_client_certificate_bound_access_tokens": true,
      //"revocation_endpoint": "https://example.com/protocol/openid-connect/revoke",
      //"revocation_endpoint_auth_methods_supported": []string{
        //"private_key_jwt",
        //"client_secret_basic",
        //"client_secret_post",
        //"tls_client_auth",
        //"client_secret_jwt",
      //},
      //"revocation_endpoint_auth_signing_alg_values_supported": []string{
        //"PS384",
        //"ES384",
        //"RS384",
        //"HS256",
        //"HS512",
        //"ES256",
        //"RS256",
        //"HS384",
        //"ES512",
        //"PS256",
        //"PS512",
        //"RS512",
      //},
      //"backchannel_logout_supported": true,
      //"backchannel_logout_session_supported": true,
      //"device_authorization_endpoint": "https://example.com/protocol/openid-connect/auth/device",
      //"backchannel_token_delivery_modes_supported": []string{
        //"poll",
        //"ping",
      //},
      //"backchannel_authentication_endpoint": "https://example.com/protocol/openid-connect/ext/ciba/auth",
      //"backchannel_authentication_request_signing_alg_values_supported": []string{
        //"PS384",
        //"ES384",
        //"RS384",
        //"ES256",
        //"RS256",
        //"ES512",
        //"PS256",
        //"PS512",
        //"RS512",
      //},
      //"require_pushed_authorization_requests": false,
      //"pushed_authorization_request_endpoint": "https://example.com/protocol/openid-connect/ext/par/request",
      //"mtls_endpoint_aliases": map[string]interface{}{
        //"token_endpoint": "https://example.com/protocol/openid-connect/token",
        //"revocation_endpoint": "https://example.com/protocol/openid-connect/revoke",
        //"introspection_endpoint": "https://example.com/protocol/openid-connect/token/introspect",
        //"device_authorization_endpoint": "https://example.com/protocol/openid-connect/auth/device",
        //"registration_endpoint": "https://example.com/clients-registrations/openid-connect",
        //"userinfo_endpoint": "https://example.com/protocol/openid-connect/userinfo",
        //"pushed_authorization_request_endpoint": "https://example.com/protocol/openid-connect/ext/par/request",
        //"backchannel_authentication_endpoint": "https://example.com/protocol/openid-connect/ext/ciba/auth",
      //},
    })
}

func tokenHandler(c echo.Context) error {
    formParams, _ := c.FormParams()

    // select random alg among available ones for JWT
    algs := []jwa.SignatureAlgorithm{
      jwa.RS256,
      jwa.ES384,
      jwa.EdDSA,
    }
    alg := algs[ rand.Intn(len(algs)) ]
    log.Printf("selected alg: %v\n", alg)

    expireDuration,_ := time.ParseDuration("1h")

    accessToken := jwt.New()
    accessToken.Set(jwt.SubjectKey, c.FormValue("code"))
    accessToken.Set(jwt.IssuerKey, BaseUrl)
    accessToken.Set(jwt.AudienceKey, `Golang Users`)
    accessToken.Set(jwt.IssuedAtKey, time.Now())
    accessToken.Set(jwt.ExpirationKey, time.Now().Add(expireDuration))
    accessToken.Set(`code`, c.FormValue("code"))

    idToken := openid.New()
    idToken.Set(jwt.SubjectKey, c.FormValue("code"))
    idToken.Set(jwt.IssuerKey, BaseUrl)
    idToken.Set(jwt.AudienceKey, `Golang Users`)
    idToken.Set(openid.NameKey, `John Doe`)
    idToken.Set(jwt.IssuedAtKey, time.Now())
    idToken.Set(jwt.ExpirationKey, time.Now().Add(expireDuration))
    idToken.Set(`code`, c.FormValue("code"))

    finalIdToken, err := CreateJWT(alg, idToken)
    if err != nil {
      log.Printf("error signing Id Token: %s\n", err)
      return err
    }

    finalAccessToken, err := CreateJWT(alg, accessToken)
    if err != nil {
      log.Printf("error signing Access Token: %s\n", err)
      return err
    }

		return c.JSON(http.StatusOK, map[string]interface{}{
      "params": formParams,
      "id_token": string(finalIdToken),
      "access_token": string(finalAccessToken),
    })
}
