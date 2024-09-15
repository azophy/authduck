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
  name               string
  requestUrl        string
  expectedStatusCode int
  expectedResponseContains   []string
  expectedResponseNotContains   []string
}

func TestServeResourceFolder(t *testing.T) {
	// Setup
	e := echo.New()

  testHandler := ServeResourceFolder("resources")

	tests := []testCase {
		{
			name:               "Valid request",
			requestUrl:        `/assets/tacit.min.css`,
			expectedStatusCode: http.StatusOK,
      expectedResponseContains: nil,
      expectedResponseNotContains: nil,
		},
	}

  runTests(tests, testHandler, e, t)
}

func TestServeResourceFile(t *testing.T) {
	e := echo.New()

  testHandler := ServeResourceFile("resources/views/home.html")

	tests := []testCase {
		{
			name:               "Valid request",
			requestUrl:        `/`,
			expectedStatusCode: http.StatusOK,
      expectedResponseContains: []string{
        "<p>To start playing with this tools, you could use any OIDC client you have",
        `<script src="/assets/htmx.min.js"></script>`, // script in layout_base.html
      },
      expectedResponseNotContains: []string{
        "random non-existing string",
      },
		},
	}

  runTests(tests, testHandler, e, t)

  testHandlerFail := ServeResourceFile("not-existing.html")

	tests = []testCase {
		{
			name:               "Requesting non-existing template",
			requestUrl:        `/`,
			expectedStatusCode: http.StatusNotFound,
      expectedResponseContains: nil,
      expectedResponseNotContains: nil,
		},
	}

  runTests(tests, testHandlerFail, e, t)
}

func TestServeResourceTemplate(t *testing.T) {
	// Setup
	e := echo.New()

	e.Renderer = SetupTemplateRegistry("resources/views/*")
  testHandler := ServeResourceTemplate("resources/views/home.html")
	e.GET("/", testHandler)

	tests := []testCase {
		{
			name:               "Valid request",
			requestUrl:        `/`,
			expectedStatusCode: http.StatusOK,
      expectedResponseContains: nil,
      expectedResponseNotContains: nil,
		},
		{
			name:               "Not found",
			requestUrl:        `/not-exists`,
			expectedStatusCode: http.StatusNotFound,
      expectedResponseContains: nil,
      expectedResponseNotContains: nil,
		},
	}

  runTests(tests, testHandler, e, t)
}

func runTests(testcases []testCase, testHandler echo.HandlerFunc, e *echo.Echo, t *testing.T) {
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.requestUrl, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := testHandler(c)
			if err != nil {
				t.Errorf("handler returned error: %v", err)
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
