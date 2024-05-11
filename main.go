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

	e.Logger.Fatal(e.Start(":3000"))
}
