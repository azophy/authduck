package main

import (
  "os"
  "log"
  "time"
  "errors"
  "strconv"
  "crypto/rsa"
  "crypto/rand"
  "crypto/ecdsa"
  "crypto/ed25519"
  "crypto/elliptic"
  "database/sql"

	"github.com/labstack/echo/v4"
  _ "github.com/mattn/go-sqlite3"
  "github.com/lestrrat-go/jwx/v2/jwa"
  "github.com/lestrrat-go/jwx/v2/jwk"
  "github.com/lestrrat-go/jwx/v2/jwt"
)

var (
  BaseUrl = GetEnvOrDefault("BASE_URL", "http://localhost:3000")
  JwkIdPrefix = GetEnvOrDefault("JWK_ID_PREFIX", "authduck")

  HistoryRepository *HistoryModel
  CodeExchangeRepository *CodeExchangeModel

  PublicJWKS = jwk.NewSet()
  RSAPrivateJWK, RSAPublicJWK, ECPrivateJWK, ECPublicJWK, EDPrivateJWK, EDPublicJWK jwk.Key

  timestamp = strconv.FormatInt(time.Now().Unix(), 10)
  RSAKeyId = JwkIdPrefix + "-rsa-key-" + timestamp
  ECKeyId = JwkIdPrefix + "-ec-key-" + timestamp
  EDKeyId = JwkIdPrefix + "-ed-key-" + timestamp

  DBFilePath = GetEnvOrDefault("DB_FILE_PATH", ":memory:")
  DBConn *sql.DB
)

func InitiateGlobalVars() error {
  // Unix() return int64, so we could not use Itoa(). Instead we used FormatInt() here

  var err error
  log.Println("generating RSA private key...")
  rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
  if err != nil {
    log.Printf("failed to generate RSA private key: %s", err)
    return err
  }

  RSAPrivateJWK, err = jwk.FromRaw(rsaPrivateKey)
  if err != nil {
    log.Printf("failed to create RSA Private JWK: %s\n", err)
    return err
  }
  RSAPrivateJWK.Set(jwk.KeyIDKey, RSAKeyId)

  RSAPublicJWK, err = RSAPrivateJWK.PublicKey()
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

  ECPrivateJWK, err = jwk.FromRaw(ecPrivateKey)
  if err != nil {
    log.Printf("failed to create EC Private JWK: %s\n", err)
    return err
  }
  ECPrivateJWK.Set(jwk.KeyIDKey, ECKeyId)

  ECPublicJWK, err = ECPrivateJWK.PublicKey()
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

  EDPrivateJWK, err = jwk.FromRaw(edPrivateKey)
  if err != nil {
    log.Printf("failed to create ED Private JWK: %s\n", err)
    return err
  }
  EDPrivateJWK.Set(jwk.KeyIDKey, EDKeyId)

  EDPublicJWK, err = jwk.FromRaw(edPublicKey)
  if err != nil {
    log.Printf("failed to create ED Public JWK: %s\n", err)
    return err
  }
  EDPublicJWK.Set(jwk.KeyIDKey, EDKeyId)

  _ = PublicJWKS.AddKey(RSAPublicJWK)
  _ = PublicJWKS.AddKey(ECPublicJWK)
  _ = PublicJWKS.AddKey(EDPublicJWK)

  log.Println("setting up db")
  DBConn, err := sql.Open("sqlite3", DBFilePath)
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
  if (alg == jwa.RS256) {
    return jwt.Sign(
      token,
      jwt.WithKey(jwa.RS256, RSAPrivateJWK),
    )
  } else if (alg == jwa.ES384) {
    return jwt.Sign(
      token,
      jwt.WithKey(jwa.ES384, ECPrivateJWK),
    )
  } else if (alg == jwa.EdDSA) {
    return jwt.Sign(
      token,
      jwt.WithKey(jwa.EdDSA, EDPrivateJWK),
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
