package main

import (
  "log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
  APP_PORT := GetEnvOrDefault("APP_PORT", "3000")
  err := InitiateGlobalVars()
  if err != nil {
    log.Printf("failed initiating global vars: %s\n", err)
    return
  }

	e := echo.New()

  e.Renderer = SetupTemplateRegistry("resources/views/*")
  RegisterHistoryHandlers(e)

	e.GET("/", ServeResourceTemplate("resources/views/home.html", nil))
	e.GET("/assets/*", ServeResourceFolder("resources"))

	e.GET("/.well-known/certs", func (c echo.Context) error {
		return c.JSON(http.StatusOK, PublicJWKS)
  })

  RegisterGenericOAuthHandlers(e)

	e.Logger.Fatal(e.Start(":" + APP_PORT))
}


