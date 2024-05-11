package main

import (
	//"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.File("/", "resources/pages/index.html")
	e.Static("/assets", "resources/assets")

	e.GET("/.well-known/openid-configuration", OpenidconfigHandler)

	//e.GET("/auth/callback", CallbackHandler)
  e.File("/auth/callback", "resources/pages/callback.html")
	e.POST("/auth/token", TokenHandler)

	e.Logger.Fatal(e.Start(":3000"))
}
