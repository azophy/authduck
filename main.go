package main

import (
  "os"
  "log"
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed resources
var embededFiles embed.FS

func main() {
  APP_PORT := GetEnvOrDefault("APP_PORT", "3000")
  GenerateSigningKeys()
	e := echo.New()

  assetHandler := http.FileServer(getFileSystem(false))
	//e.GET("/", echo.WrapHandler(http.StripPrefix("/pages/", assetHandler)))
	e.GET("/", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.ServeFileFS(w, r, embededFiles, "resources/pages/index.html")
  })))
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	//e.File("/", "resources/pages/index.html")
	//e.Static("/assets", "resources/assets")

	e.GET("/.well-known/certs", func (c echo.Context) error {
		return c.JSON(http.StatusOK, PublicJWKS)
  })

  GeneralOAuthModuleInit(e)

	e.Logger.Fatal(e.Start(":" + APP_PORT))
}

func getFileSystem(useOS bool) http.FileSystem {
	if useOS {
		log.Print("using live mode")
		return http.FS(os.DirFS("resources"))
	}

	log.Print("using embed mode")
	fsys, err := fs.Sub(embededFiles, "resources")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
