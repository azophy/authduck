package main

import (
  "log"
  "time"
  "net/http"
  "math/rand"

	"github.com/labstack/echo/v4"
  "github.com/lestrrat-go/jwx/v2/jwa"
  "github.com/lestrrat-go/jwx/v2/jwt"
  "github.com/lestrrat-go/jwx/v2/jwt/openid"
)

func TokenHandler(c echo.Context) error {
    formParams, _ := c.FormParams()
    algs := []jwa.SignatureAlgorithm{
      jwa.RS256,
      jwa.ES384,
      jwa.EdDSA,
    }
    //randomIndex := rand.Intn(len(algs))
    //log.Printf("random index: %v\n", randomIndex)
    //alg := algs[2]
    alg := algs[ rand.Intn(len(algs)) ]
    log.Printf("selected alg: %v\n", alg)

    expireDuration,_ := time.ParseDuration("1h")

    accessToken := jwt.New()
    accessToken.Set(jwt.SubjectKey, `https://github.com/azophy/authduck`)
    accessToken.Set(jwt.AudienceKey, `Golang Users`)
    accessToken.Set(jwt.IssuedAtKey, time.Now())
    accessToken.Set(jwt.ExpirationKey, time.Now().Add(expireDuration))
    accessToken.Set(`code`, c.FormValue("code"))

    idToken := openid.New()
    idToken.Set(jwt.SubjectKey, `https://github.com/azophy/authduck`)
    idToken.Set(jwt.AudienceKey, `Golang Users`)
    idToken.Set(openid.NameKey, `John Doe`)
    idToken.Set(jwt.IssuedAtKey, time.Now())
    idToken.Set(jwt.ExpirationKey, time.Now().Add(expireDuration))
    idToken.Set(`code`, c.FormValue("code"))

    finalIdToken, err := CreateJWT(alg, idToken)
    if err != nil {
      log.Printf("error signing Id Token: %s\n", err)
      return err
    }

    finalAccessToken, err := CreateJWT(alg, accessToken)
    if err != nil {
      log.Printf("error signing Access Token: %s\n", err)
      return err
    }

		return c.JSON(http.StatusOK, map[string]interface{}{
      "params": formParams,
      "id_token": string(finalIdToken),
      "access_token": string(finalAccessToken),
    })
}
