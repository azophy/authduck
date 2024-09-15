package main

import (
	//"database/sql"
	//"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	//"github.com/stretchr/testify/assert"
)

type testCase struct {
  name               string
  requestUrl        string
  //requestBody        string
  expectedStatusCode int
  //expectedResponse   string
}


func TestServeResourceFolder(t *testing.T) {
	// Setup
	e := echo.New()

  testHandler := ServeResourceFolder("resources")

	tests := []testCase {
		{
			name:               "Valid request",
			requestUrl:        `/assets/tacit.min.css`,
			//requestBody:        `/assets/tacit.min.css`,
			expectedStatusCode: http.StatusOK,
			//expectedResponse:   `{"id_token":"valid_id_token","access_token":"valid_access_token"}`,
		},
		{
			name:               "Not found",
			requestUrl:        `/random-non-existing-file.min.css`,
			expectedStatusCode: http.StatusNotFound,
			//expectedResponse:   `{"id_token":"valid_id_token","access_token":"valid_access_token"}`,
		},
	}

  runTests(tests, testHandler, e, t)
}

func runTests(testcases []testCase, testHandler echo.HandlerFunc, e *echo.Echo, t *testing.T) {
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.requestUrl, nil)
			//req := httptest.NewRequest(http.MethodGet, tt.requestUrl, strings.NewReader(tt.requestBody))
			//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := testHandler(c)
			if err != nil {
				t.Errorf("tokenHandler() returned error: %v", err)
			}

			if rec.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, rec.Code)
			}

			//if strings.TrimSpace(rec.Body.String()) != tt.expectedResponse {
				//t.Errorf("expected response %q, got %q", tt.expectedResponse, rec.Body.String())
			//}
		})
	}
}
