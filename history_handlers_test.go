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

func setupTestDB(t *testing.T) *sql.DB {
	err := InitiateDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to initiate test database: %v", err)
	}

	return DBConn
}

func TestHistoryDetailHandler(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = SetupTemplateRegistry("resources/views/*")
	req := httptest.NewRequest(http.MethodGet, "/manage/history?id=testClient&from=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create a real HistoryModel instance
	historyModel, err := NewHistoryModel(db)
	assert.NoError(t, err)

	// Insert test data using the Record method
	err = historyModel.Record("testClient", "GET", "/test", "", "", "")
	assert.NoError(t, err)

	// Replace the global HistoryRepository with the test instance
	HistoryRepository = historyModel

	// Test
	if assert.NoError(t, historyDetailHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// You can add more specific assertions here about the response body
		// For example, you could parse the JSON response and check its contents
	}
}

func TestExtractBasicAuthMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test with valid Basic Auth
	username, password := "testUser", "testPass"
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Set("Authorization", "Basic "+auth)

	handler := extractBasicAuthMiddleware(func(c echo.Context) error {
		assert.Equal(t, username, c.Get("basic-auth-username"))
		assert.Equal(t, password, c.Get("basic-auth-password"))
		return nil
	})

	assert.NoError(t, handler(c))
}

func TestHistoryRecorderMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/?client_id=testClient", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create a real HistoryModel instance
	historyModel, err := NewHistoryModel(db)
	assert.NoError(t, err)

	// Replace the global HistoryRepository with the test instance
	HistoryRepository = historyModel

	// Test
	handler := historyRecorderMiddleware(func(c echo.Context) error {
		return nil
	})

	assert.NoError(t, handler(c))

	// Verify that the history was recorded
	histories, err := HistoryRepository.All("testClient", "0")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(histories))
	assert.Equal(t, "testClient", histories[0].ClientId)
	assert.Equal(t, "POST", histories[0].HTTPMethod)
	assert.Equal(t, "/?client_id=testClient", histories[0].Url)
}
