package main

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)


func TestServeResourceFolder(t *testing.T) {
	// Setup
	e := echo.New()

  testHandler := ServeResourceFolder("resources")

	req := httptest.NewRequest(http.MethodGet, "/assets/tacit.min.css", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	if assert.NoError(t, testHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// You can add more specific assertions here about the response body
		// For example, you could parse the JSON response and check its contents
	}
}
