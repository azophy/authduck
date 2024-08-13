package main

import (
  "log"
	"net/http"
  "strconv"

  "golang.org/x/time/rate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
  APP_PORT := GetEnvOrDefault("APP_PORT", "3000")
  RATE_LIMIT := GetEnvOrDefault("RATE_LIMIT", "20")

  err := InitiateGlobalVars()
  if err != nil {
    log.Printf("failed initiating global vars: %s\n", err)
    return
  }

	e := echo.New()

  rateLimitNumber, err := strconv.Atoi(RATE_LIMIT)
  if err != nil {
    log.Printf("failed initiating rate limiter: %s\n", err)
    return
  }
  e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(rateLimitNumber))))

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


