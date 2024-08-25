package main

import (
  "log"
	"net/http"

  "golang.org/x/time/rate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
  err := InitiateGlobalVars()
  if err != nil {
    log.Printf("failed initiating global vars: %s\n", err)
    return
  }

	e := echo.New()

  e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(Config.RateLimit))))

  e.Renderer = SetupTemplateRegistry("resources/views/*")
  RegisterHistoryHandlers(e)

	e.GET("/", ServeResourceTemplate("resources/views/home.html", nil))
	e.GET("/assets/*", ServeResourceFolder("resources"))

	e.GET("/.well-known/certs", func (c echo.Context) error {
		return c.JSON(http.StatusOK, Config.PublicJWKS)
  })

  RegisterGenericOAuthHandlers(e)

	e.Logger.Fatal(e.Start(":" + Config.AppPort))
}


