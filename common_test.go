package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
)

func TestConfigInit(t *testing.T) {
	config := &ConfigType{}
	err := config.Init()
	assert.NoError(t, err)

	// Test default values
	assert.Equal(t, "3000", config.AppPort)
	assert.Equal(t, "http://localhost:3000", config.BaseUrl)
	assert.Equal(t, 20, config.RateLimit)
	assert.Equal(t, ":memory:", config.DBFilePath)
	assert.False(t, config.IsAppDebug)

	// Test JWK generation
	assert.NotNil(t, config.RSAPrivateJWK)
	assert.NotNil(t, config.RSAPublicJWK)
	assert.NotNil(t, config.ECPrivateJWK)
	assert.NotNil(t, config.ECPublicJWK)
	assert.NotNil(t, config.EDPrivateJWK)
	assert.NotNil(t, config.EDPublicJWK)
	assert.NotNil(t, config.PublicJWKS)
}

func TestGetEnvOrDefault(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	assert.Equal(t, "test_value", GetEnvOrDefault("TEST_VAR", "default"))

	// Test with non-existing environment variable
	assert.Equal(t, "default", GetEnvOrDefault("NON_EXISTING_VAR", "default"))
}

func TestIsReqFromHTMX(t *testing.T) {
	e := echo.New()

	// Test with HTMX header
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("HX-REQUEST", "true")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.True(t, IsReqFromHTMX(c))

	// Test without HTMX header
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	assert.False(t, IsReqFromHTMX(c))
}

func TestCreateJWT(t *testing.T) {
	Config = ConfigType{}
	err := Config.Init()
	assert.NoError(t, err)

	token := jwt.New()
	token.Set(jwt.SubjectKey, "test")

	// Test RS256
	signed, err := CreateJWT(jwa.RS256, token)
	assert.NoError(t, err)
	assert.NotNil(t, signed)

	// Test ES384
	signed, err = CreateJWT(jwa.ES384, token)
	assert.NoError(t, err)
	assert.NotNil(t, signed)

	// Test EdDSA
	signed, err = CreateJWT(jwa.EdDSA, token)
	assert.NoError(t, err)
	assert.NotNil(t, signed)

	// Test invalid algorithm
	_, err = CreateJWT(jwa.HS256, token)
	assert.Error(t, err)
}

func TestCreateJwtFromJson(t *testing.T) {
	Config = ConfigType{}
	err := Config.Init()
	assert.NoError(t, err)

	validPayload := `{"alg":"RS256","sub":"test"}`
	signed, err := CreateJwtFromJson(validPayload)
	assert.NoError(t, err)
	assert.NotNil(t, signed)

	invalidPayload := `{"sub":"test"}`
	_, err = CreateJwtFromJson(invalidPayload)
	assert.Error(t, err)
}

func TestIsItemInList(t *testing.T) {
	list := []string{"apple", "banana", "cherry"}

	assert.True(t, IsItemInList("banana", list))
	assert.False(t, IsItemInList("grape", list))

	numList := []int{1, 2, 3, 4, 5}
	assert.True(t, IsItemInList(3, numList))
	assert.False(t, IsItemInList(6, numList))
}

func TestInitiateDatabase(t *testing.T) {
	err := InitiateDatabase(":memory:")
	assert.NoError(t, err)
	assert.NotNil(t, DBConn)
	assert.NotNil(t, HistoryRepository)
	assert.NotNil(t, CodeExchangeRepository)
}

func TestInitiateGlobalVars(t *testing.T) {
	err := InitiateGlobalVars()
	assert.NoError(t, err)
	assert.NotNil(t, Config)
	assert.NotNil(t, DBConn)
	assert.NotNil(t, HistoryRepository)
	assert.NotNil(t, CodeExchangeRepository)
}
