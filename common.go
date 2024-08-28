package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	_ "modernc.org/sqlite"
)

type ConfigType struct {
	AppPort      string
	BaseUrl      string
	RateLimit    int
	CustomScript template.HTML
	CorsOrigins  []string
	IsAppDebug   bool

	DBFilePath string

	JwkIdPrefix                                                                       string
	RSAKeyId, ECKeyId, EDKeyId                                                        string
	RSAPrivateJWK, RSAPublicJWK, ECPrivateJWK, ECPublicJWK, EDPrivateJWK, EDPublicJWK jwk.Key
	PublicJWKS                                                                        jwk.Set
}

var (
	Config ConfigType

	DBConn                 *sql.DB
	HistoryRepository      *HistoryModel
	CodeExchangeRepository *CodeExchangeModel
)

func (ct *ConfigType) Init() error {
	var err error
	ct.AppPort = GetEnvOrDefault("APP_PORT", "3000")
	ct.BaseUrl = GetEnvOrDefault("BASE_URL", "http://localhost:3000")
	ct.CustomScript = template.HTML(GetEnvOrDefault("CUSTOM_SCRIPT", ""))
	ct.DBFilePath = GetEnvOrDefault("DB_FILE_PATH", ":memory:")
	ct.IsAppDebug = (GetEnvOrDefault("APP_DEBUG", "false") == "true")

	RateLimitEnv := GetEnvOrDefault("RATE_LIMIT", "20")
	ct.RateLimit, err = strconv.Atoi(RateLimitEnv)
	if err != nil {
		log.Printf("failed parsing rate limiter config: %s\n", err)
		return err
	}

	ct.CorsOrigins = strings.Split(GetEnvOrDefault("CORS_ORIGINS", ""), ";")

	ct.JwkIdPrefix = GetEnvOrDefault("JWK_ID_PREFIX", "authduck")

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	ct.RSAKeyId = ct.JwkIdPrefix + "-rsa-key-" + timestamp
	ct.ECKeyId = ct.JwkIdPrefix + "-ec-key-" + timestamp
	ct.EDKeyId = ct.JwkIdPrefix + "-ed-key-" + timestamp
	ct.PublicJWKS = jwk.NewSet()

	log.Println("generating RSA private key...")
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("failed to generate RSA private key: %s", err)
		return err
	}

	ct.RSAPrivateJWK, err = jwk.FromRaw(rsaPrivateKey)
	if err != nil {
		log.Printf("failed to create RSA Private JWK: %s\n", err)
		return err
	}
	ct.RSAPrivateJWK.Set(jwk.KeyIDKey, ct.RSAKeyId)

	ct.RSAPublicJWK, err = ct.RSAPrivateJWK.PublicKey()
	if err != nil {
		log.Printf("failed to create RSA Public JWK: %s\n", err)
		return err
	}

	log.Println("generating EC private key...")
	ecPrivateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		log.Printf("failed to generate new ECDSA private key: %s\n", err)
		return err
	}

	ct.ECPrivateJWK, err = jwk.FromRaw(ecPrivateKey)
	if err != nil {
		log.Printf("failed to create EC Private JWK: %s\n", err)
		return err
	}
	ct.ECPrivateJWK.Set(jwk.KeyIDKey, ct.ECKeyId)

	ct.ECPublicJWK, err = ct.ECPrivateJWK.PublicKey()
	if err != nil {
		log.Printf("failed to create EC Public JWK: %s\n", err)
		return err
	}

	log.Println("Generating ED25519 keypair...")
	edPublicKey, edPrivateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Printf("failed to generate new ED25519 keypair: %s\n", err)
		return err
	}

	ct.EDPrivateJWK, err = jwk.FromRaw(edPrivateKey)
	if err != nil {
		log.Printf("failed to create ED Private JWK: %s\n", err)
		return err
	}
	ct.EDPrivateJWK.Set(jwk.KeyIDKey, ct.EDKeyId)

	ct.EDPublicJWK, err = jwk.FromRaw(edPublicKey)
	if err != nil {
		log.Printf("failed to create ED Public JWK: %s\n", err)
		return err
	}
	ct.EDPublicJWK.Set(jwk.KeyIDKey, ct.EDKeyId)

	_ = ct.PublicJWKS.AddKey(ct.RSAPublicJWK)
	_ = ct.PublicJWKS.AddKey(ct.ECPublicJWK)
	_ = ct.PublicJWKS.AddKey(ct.EDPublicJWK)
	return nil
}

func (ct *ConfigType) GetCORSConfig() middleware.CORSConfig {
	allowedMethods := []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete}

	log.Printf("origij lemgthv %v", len(ct.CorsOrigins))

	if len(ct.CorsOrigins) > 0 {
		return middleware.CORSConfig{
			AllowOrigins: ct.CorsOrigins,
			AllowMethods: allowedMethods,
		}
	} else {
		return middleware.CORSConfig{
			AllowMethods: allowedMethods,
		}
	}
}

func InitiateGlobalVars() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %s", err)
		log.Println("Continuing without env file...")
	}

	// initiate global config
	Config.Init()

	log.Println("setting up db")
	DBConn, err := sql.Open("sqlite", Config.DBFilePath)
	if err != nil {
		return err
	}

	HistoryRepository, err = NewHistoryModel(DBConn)
	if err != nil {
		log.Printf("failed to initiate HistoryModel: %s\n", err)
		return err
	}

	CodeExchangeRepository, err = NewCodeExchangeModel(DBConn)
	if err != nil {
		log.Printf("failed to initiate CodeExchangeModel: %s\n", err)
		return err
	}

	return nil
}

func GetEnvOrDefault(varName string, defaultValue string) string {
	if os.Getenv(varName) != "" {
		return os.Getenv(varName)
	} else {
		return defaultValue
	}
}

func IsReqFromHTMX(c echo.Context) bool {
	htmxHeader := c.Request().Header.Get("HX-REQUEST")
	return (htmxHeader != "")
}

func CreateJWT(alg jwa.SignatureAlgorithm, token jwt.Token) ([]byte, error) {
	if alg == jwa.RS256 {
		return jwt.Sign(
			token,
			jwt.WithKey(jwa.RS256, Config.RSAPrivateJWK),
		)
	} else if alg == jwa.ES384 {
		return jwt.Sign(
			token,
			jwt.WithKey(jwa.ES384, Config.ECPrivateJWK),
		)
	} else if alg == jwa.EdDSA {
		return jwt.Sign(
			token,
			jwt.WithKey(jwa.EdDSA, Config.EDPrivateJWK),
		)
	}

	return nil, errors.New("Invalid jwt signature")
}

func CreateJwtFromJson(rawPayload string) ([]byte, error) {
	token, err := jwt.Parse([]byte(rawPayload), jwt.WithVerify(false))
	if err != nil {
		log.Printf("error on generating JWT: %v", err)
		return nil, err
	}
	alg, ok := token.Get("alg")
	if !ok {
		return nil, errors.New("no alg defined inside JSON payload")
	}
	return CreateJWT(jwa.SignatureAlgorithm(alg.(string)), token)
}
