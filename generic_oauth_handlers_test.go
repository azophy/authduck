package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTokenHandler(t *testing.T) {
	e := echo.New()
	err := InitiateDatabase(":memory:")
	if err != nil {
		t.Errorf("encounter error: %v", err)
	}
	err = CodeExchangeRepository.Add(
		"validClientId",
		"validCode",
		`{"id_token":"valid_id_token","access_token":"valid_access_token"}`,
	)
	if err != nil {
		t.Errorf("error on generating exchange payload: %v", err)
	}

	tests := []struct {
		name               string
		requestBody        string
		expectedStatusCode int
		expectedResponse   string
	}{
		// negative states
		{
			name:               "Invalid grant type",
			requestBody:        `{"client_id":"validClientId", "code":"validCode","grant_type":"invalid_grant_type"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"bad request","error_description":"invalid grant type"}`,
		},
		{
			name:               "Invalid JSON request",
			requestBody:        `invalid json`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "bad request",
		},

		// authorization code flow
		{
			name:               "Authorization Code Flow: Valid request",
			requestBody:        `{"client_id":"validClientId", "code":"validCode","grant_type":"authorization_code"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id_token":"valid_id_token","access_token":"valid_access_token"}`,
		},
		{
			name:               "Authorization Code Flow: Client Id Not found",
			requestBody:        `{"client_id":"invalidClientId", "code":"someCode","grant_type":"authorization_code"}`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"error":"not found","error_description":"client id-secret pair not found"}`,
		},
		{
			name:               "Authorization Code Flow: Code Not found",
			requestBody:        `{"client_id":"validClientId", "code":"someCode","grant_type":"authorization_code"}`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"error":"not found","error_description":"client id-secret pair not found"}`,
		},

		// client credential flow
		{
			name:               "Client Credential Flow: Client id registered from other flows",
			requestBody:        `{"grant_type":"client_credentials","client_id":"validClientId", "client_secret":"someSecret"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id_token":"valid_id_token","access_token":"valid_access_token"}`,
		},
		{
			name:               "Client Credential Flow: Unregistered Client Id should also works",
			requestBody:        `{"grant_type":"client_credentials","client_id":"someClientId", "client_secret":"someSecret"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id_token":"valid_id_token","access_token":"valid_access_token"}`,
		},
		{
			name:               "Client Credential Flow: Client Id is required",
			requestBody:        `{"grant_type":"client_credentials"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"not found","error_description":"client id field is empty"}`,
		},
		{
			name:               "Client Credential Flow: Client Secret is required",
			requestBody:        `{"grant_type":"client_credentials","client_id":"someClientId"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"not found","error_description":"client secret field is empty"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := tokenHandler(c)
			if err != nil {
				t.Errorf("tokenHandler() returned error: %v", err)
			}

			if rec.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, rec.Code)
			}

			if strings.TrimSpace(rec.Body.String()) != tt.expectedResponse {
				t.Errorf("expected response %q, got %q", tt.expectedResponse, rec.Body.String())
			}
		})
	}
}

func TestGetOpenidConfig(t *testing.T) {
	// Set up test configuration
	Config.BaseUrl = "https://example.com"

	config := getOpenidConfig()

	assert.Equal(t, "https://example.com", config["issuer"])
	assert.Equal(t, "https://example.com/case/generic/auth/callback", config["authorization_endpoint"])
	assert.Equal(t, "https://example.com/case/generic/auth/token", config["token_endpoint"])
	assert.Equal(t, SUPPORTED_GRANT_TYPES, config["grant_types_supported"])
	assert.Equal(t, SUPPORTED_RESPONSE_TYPES, config["response_types_supported"])
	assert.Equal(t, SUPPORTED_SIGNING_ALGS, config["id_token_signing_alg_values_supported"])
}

func TestOpenidconfigHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/.well-known/openid-configuration", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, openidconfigHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, Config.BaseUrl, response["issuer"])
		assert.Equal(t, Config.BaseUrl+"/case/generic/auth/callback", response["authorization_endpoint"])
		assert.Equal(t, Config.BaseUrl+"/case/generic/auth/token", response["token_endpoint"])
	}
}

func TestCallbackPostHandler(t *testing.T) {
	e := echo.New()
	err := InitiateDatabase(":memory:")
	assert.NoError(t, err)

	err = Config.Init()
	assert.NoError(t, err, "Failed to initialize Config")

	callbackPayload := `{"code":"test_code","redirect_uri":"https://client.example.com/callback"}`
	accessTokenPayload := `{"sub":"user123","exp":1735689600,"alg":"RS256"}`
	idTokenPayload := `{"sub":"user123","aud":"client123","exp":1735689600,"alg":"RS256"}`

	form := url.Values{}
	form.Set("callback_payload", callbackPayload)
	form.Set("client_id", "test_client_id")
	form.Set("access_token_payload", accessTokenPayload)
	form.Set("id_token_payload", idTokenPayload)

	req := httptest.NewRequest(http.MethodPost, "/auth/callback", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err = callbackPostHandler(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	location := rec.Header().Get("Location")
	assert.Contains(t, location, "https://client.example.com/callback")
	assert.Contains(t, location, "code=test_code")

	// Verify that the code exchange payload was stored
	exchange, err := CodeExchangeRepository.GetOne("test_client_id", "test_code")
	assert.NoError(t, err)
	assert.NotNil(t, exchange)

	var tokenResponse map[string]string
	err = json.Unmarshal([]byte(exchange.Payload), &tokenResponse)
	assert.NoError(t, err)
	assert.Contains(t, tokenResponse, "id_token")
	assert.Contains(t, tokenResponse, "access_token")
}
