package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

const (
	routeParent = "/case/generic"
)

func RegisterGenericOAuthHandlers(app *echo.Echo) {
	e := app.Group(routeParent)

	e.GET("/.well-known/openid-configuration", openidconfigHandler)
	e.GET("/auth/callback", ServeResourceTemplate("resources/views/generic_callback.html"))
	e.POST("/auth/callback", callbackPostHandler)
	e.POST("/auth/token", tokenHandler)
}

func openidconfigHandler(c echo.Context) error {
	BaseUrl := Config.BaseUrl
	return c.JSON(http.StatusOK, echo.Map{
		"issuer":                 BaseUrl,
		"authorization_endpoint": BaseUrl + routeParent + "/auth/callback",
		"token_endpoint":         BaseUrl + routeParent + "/auth/token",
		"introspection_endpoint": BaseUrl + routeParent + "/auth/introspect",
		"userinfo_endpoint":      BaseUrl + routeParent + "/auth/userinfo",
		"end_session_endpoint":   BaseUrl + routeParent + "/auth/logout",
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
		//"mtls_endpoint_aliases": echo.Map{
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

func callbackPostHandler(c echo.Context) error {
	callbackPayloadRaw := c.FormValue("callback_payload")
	clientId := c.FormValue("client_id")

	// prepare callback
	var callbackPayload echo.Map
	_ = json.Unmarshal([]byte(callbackPayloadRaw), &callbackPayload)

	q := url.Values{}
	for k, v := range callbackPayload {
		q.Set(k, v.(string))
	}

	redirectUrl, err := url.Parse(callbackPayload["redirect_uri"].(string))
	if err != nil {
		log.Printf("error on redirect_url: %v", err)
		return err
	}
	redirectUrl.RawQuery = q.Encode()

	// store code exchange payload
	accessTokenPayloadRaw := c.FormValue("access_token_payload")
	accessToken, err := CreateJwtFromJson(accessTokenPayloadRaw)
	if err != nil {
		log.Printf("error on generating access token: %v", err)
		return err
	}
	idTokenPayloadRaw := c.FormValue("id_token_payload")
	idToken, err := CreateJwtFromJson(idTokenPayloadRaw)
	if err != nil {
		log.Printf("error on generating id token: %v", err)
		return err
	}

	codeExchangePayload, _ := json.Marshal(echo.Map{
		"id_token":     string(idToken),
		"access_token": string(accessToken),
	})
	err = CodeExchangeRepository.Add(
		clientId,
		callbackPayload["code"].(string),
		string(codeExchangePayload),
	)
	if err != nil {
		log.Printf("error on generating exchange payload: %v", err)
		return err
	}

	return c.Redirect(http.StatusSeeOther, redirectUrl.String())
}

type tokenRequest struct {
	ClientId     string `json:"client_id" form:"client_id" query:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret" query:"client_secret"`
	Code         string `json:"code" form:"code" query:"code"`
	GrantType    string `json:"grant_type" form:"grant_type" query:"grant_type"`
	Audience     string `json:"audience" form:"audience" query:"audience"`
}

func tokenHandler(c echo.Context) error {
	var payload tokenRequest
	err := c.Bind(&payload)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	log.Printf("payload: %v", payload)

	codeExchange, err := CodeExchangeRepository.GetOne(payload.ClientId, payload.Code)
	if err != nil {
		log.Printf("error on retrieving exchange payload: %v", err)

		if errors.Is(err, ErrNotExists) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error":             "not found",
				"error_description": "client id-secret pair not found",
			})
		} else {
			return err
		}
	}

	return c.JSONBlob(http.StatusOK, []byte(codeExchange.Payload))
}
