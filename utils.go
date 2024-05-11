package main

import (
  "os"
  "log"
  "errors"
  "crypto/rsa"
  "crypto/rand"
  "crypto/ecdsa"
  "crypto/ed25519"
  "crypto/elliptic"

  "github.com/lestrrat-go/jwx/v2/jwa"
  "github.com/lestrrat-go/jwx/v2/jwk"
  "github.com/lestrrat-go/jwx/v2/jwt"
)

var (
  BaseUrl = GetEnvOrDefault("BASE_URL", "http://localhost:3000")
  PublicJWKS = jwk.NewSet()

  RSAPrivateKey *rsa.PrivateKey
  ECPrivateKey *ecdsa.PrivateKey
  EDPrivateKey ed25519.PrivateKey
)

func GenerateSigningKeys() {
  var err error
  log.Println("generating RSA private key...")
  RSAPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
  if err != nil {
    log.Printf("failed to generate private key: %s", err)
    return
  }
  RSAPublicKey := RSAPrivateKey.Public()
  RSAPublicJWK, err := jwk.FromRaw(RSAPublicKey)
  if err != nil {
    log.Printf("failed to create RSA Public JWK: %s\n", err)
    return
  }

  log.Println("generating EC private key...")
  ECPrivateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
  if err != nil {
    log.Printf("failed to generate new ECDSA private key: %s\n", err)
    return
  }

  ECPublicKey := ECPrivateKey.Public()

  ECPublicJWK, err := jwk.FromRaw(ECPublicKey)
  if err != nil {
    log.Printf("failed to create EC Public JWK: %s\n", err)
    return
  }

  log.Println("Generating ED25519 keypair...")
  var EDPublicKey ed25519.PublicKey
  EDPublicKey, EDPrivateKey, err = ed25519.GenerateKey(rand.Reader)
  if err != nil {
    log.Printf("failed to generate new ED25519 keypair: %s\n", err)
    return
  }

  EDPublicJWK, err := jwk.FromRaw(EDPublicKey)
  if err != nil {
    log.Printf("failed to create ED Public JWK: %s\n", err)
    return
  }

  _ = PublicJWKS.AddKey(RSAPublicJWK)
  _ = PublicJWKS.AddKey(ECPublicJWK)
  _ = PublicJWKS.AddKey(EDPublicJWK)
}

func CreateJWT(alg jwa.SignatureAlgorithm, token jwt.Token) ([]byte, error) {
  if (alg == jwa.RS256) {
    return jwt.Sign(token, jwt.WithKey(jwa.RS256, RSAPrivateKey))
  } else if (alg == jwa.ES384) {
    return jwt.Sign(token, jwt.WithKey(jwa.ES384, ECPrivateKey))
  } else if (alg == jwa.EdDSA) {
    return jwt.Sign(token, jwt.WithKey(jwa.EdDSA, EDPrivateKey))
  }

  return nil, errors.New("Invalid jwt signature")
}

func GetEnvOrDefault(varName string, defaultValue string) string {
  if os.Getenv(varName) != "" {
    return os.Getenv(varName)
  } else {
    return defaultValue
  }
}
