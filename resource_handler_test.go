package main

import (
	//"database/sql"
	//"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name                        string
	handlerFunc                 echo.HandlerFunc
	requestUrl                  string
	expectedStatusCode          int
	expectedResponseContains    []string
	expectedResponseNotContains []string
	expectedErrorString         string
}

func TestServeResourceFolder(t *testing.T) {
	// Setup
	e := echo.New()

	testHandler := ServeResourceFolder("resources")

	tests := []testCase{
		{
			name:                        "Valid request",
			handlerFunc:                 testHandler,
			requestUrl:                  "/assets/tacit.min.css",
			expectedStatusCode:          http.StatusOK,
			expectedResponseContains:    nil,
			expectedResponseNotContains: nil,
			expectedErrorString:         "",
		},
		{
			name:                        "Not Found",
			handlerFunc:                 testHandler,
			requestUrl:                  "/non-existing-url",
			expectedStatusCode:          http.StatusNotFound,
			expectedResponseContains:    nil,
			expectedResponseNotContains: nil,
			expectedErrorString:         "",
		},
	}

	runTests(tests, e, t)
}

func TestServeResourceTemplate(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = SetupTemplateRegistry("resources/views/*")

	tests := []testCase{
		{
			name:               "Valid full template",
			handlerFunc:        ServeResourceTemplate("resources/views/home.html"),
			requestUrl:         "/",
			expectedStatusCode: http.StatusOK,
			expectedResponseContains: []string{
				"<p>To start playing with this tools, you could use any OIDC client you have",
				`<script src="/assets/htmx.min.js"></script>`, // script in layout_base.html
			},
			expectedResponseNotContains: []string{
				"random non-existing string",
			},
			expectedErrorString: "",
		},
		{
			name:                        "Not Existing Template",
			handlerFunc:                 ServeResourceTemplate("non-existing.html"),
			requestUrl:                  "/",
			expectedStatusCode:          http.StatusOK,
			expectedResponseContains:    nil,
			expectedResponseNotContains: nil,
			expectedErrorString:         "Template not found -> non-existing.html",
		},
		{
			name:               "Partial Template",
			handlerFunc:        ServeResourceTemplate("resources/views/home.html#partial"),
			requestUrl:         "/",
			expectedStatusCode: http.StatusOK,
			expectedResponseContains: []string{
				"<p>To start playing with this tools, you could use any OIDC client you have",
			},
			expectedResponseNotContains: []string{
				`<script src="/assets/htmx.min.js"></script>`, // script in layout_base.html
				"random non-existing string",
			},
			expectedErrorString: "",
		},
	}

	runTests(tests, e, t)
}

func runTests(testcases []testCase, e *echo.Echo, t *testing.T) {
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.requestUrl, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := tt.handlerFunc(c)
			if tt.expectedErrorString != "" {
				assert.EqualError(t, err, tt.expectedErrorString)
			} else {
				assert.NoErrorf(t, err, "handler returned error: %v", err)
			}

			if rec.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, rec.Code)
			}

			if tt.expectedResponseContains != nil {
				for _, contains := range tt.expectedResponseContains {
					assert.Contains(t, rec.Body.String(), contains)
				}
			}

			if tt.expectedResponseNotContains != nil {
				for _, contains := range tt.expectedResponseNotContains {
					assert.NotContains(t, rec.Body.String(), contains)
				}
			}
		})
	}
}
