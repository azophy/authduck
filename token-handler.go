package main

import (
  "log"
  //"time"
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
    alg := algs[rand.Intn(len(algs))]

    accessToken := jwt.New()
    accessToken.Set(jwt.SubjectKey, `https://github.com/lestrrat-go/jwx/v2/jwt`)
    accessToken.Set(jwt.AudienceKey, `Golang Users`)
    //accessToken.Set(jwt.IssuedAtKey, time.Unix(aLongLongTimeAgo, 0))
    accessToken.Set(`code`, c.FormValue("code"))

    idToken := openid.New()
    idToken.Set(jwt.SubjectKey, `https://github.com/lestrrat-go/jwx/v2/jwt`)
    idToken.Set(jwt.AudienceKey, `Golang Users`)
    idToken.Set(openid.NameKey, `John Doe`)
    //idToken.Set(jwt.IssuedAtKey, time.Unix(aLongLongTimeAgo, 0))
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
      "id_token": finalIdToken,
      "access_token": finalAccessToken,
    })
}
