package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
  APP_PORT := GetEnvOrDefault("APP_PORT", "3000")
  GenerateSigningKeys()

	e := echo.New()
	e.File("/", "resources/pages/index.html")
	e.Static("/assets", "resources/assets")

	e.GET("/.well-known/certs", func (c echo.Context) error {
		return c.JSON(http.StatusOK, PublicJWKS)
  })

  GeneralOAuthModuleInit(e)

	e.Logger.Fatal(e.Start(":" + APP_PORT))
}
