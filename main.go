package main

import (
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

  assetHandler := http.FileServer(getFileSystem())
	e.GET("/", serveEmbededFile("resources/pages/index.html"))
	e.GET("/assets/*", echo.WrapHandler(assetHandler))

	e.GET("/.well-known/certs", func (c echo.Context) error {
		return c.JSON(http.StatusOK, PublicJWKS)
  })

  GeneralOAuthModuleInit(e)

	e.Logger.Fatal(e.Start(":" + APP_PORT))
}

func serveEmbededFile(path string) echo.HandlerFunc {
  return echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.ServeFileFS(w, r, embededFiles, path)
  }))
}

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embededFiles, "resources")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
