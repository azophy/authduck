package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
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
		{
			name:               "Valid request",
			requestBody:        `{"client_id":"validClientId", "code":"validCode"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id_token":"valid_id_token","access_token":"valid_access_token"}`,
		},
		{
			name:               "Invalid JSON request",
			requestBody:        `invalid json`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "bad request",
		},
		{
			name:               "Client Id Not found",
			requestBody:        `{"client_id":"invalidClientId", "code":"someCode"}`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"error":"not found","error_description":"client id-secret pair not found"}`,
		},
		{
			name:               "Code Not found",
			requestBody:        `{"client_id":"validClientId", "code":"someCode"}`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"error":"not found","error_description":"client id-secret pair not found"}`,
		},
		{
			name:               "Internal server error",
			requestBody:        `{"client_id":"someClientId", "code":"someCode"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "internal server error",
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
