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

  e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
    c.Set("basic-auth-username", username)
    c.Set("basic-auth-password", password)
    return true, nil
  }))
  e.Use(HistoryRecorderMiddleware)

  assetHandler := http.FileServer(getFileSystem())
	e.GET("/", serveEmbededFile("resources/pages/index.html"))
	e.GET("/assets/*", echo.WrapHandler(assetHandler))

	e.GET("/.well-known/certs", func (c echo.Context) error {
		return c.JSON(http.StatusOK, PublicJWKS)
  })

  e.GET("/manage/history/:client_id/:from", func (c echo.Context) error {
    clientId := c.Param("client_id")
    from := c.Param("from")
    return c.JSON(http.StatusOK, HistoryRepository.All(clientId, from))
  })
  RegisterGeneralOAuthModule(e)

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

func HistoryRecorderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
    var clientId string

    if v := c.Get("basic-auth-username"); v != "" {
      clientId = v
    }
    if v := c.QueryParam("client_id"); v != "" {
      clientId = v
    }
    if v := c.FormValue("client_id"); v != "" {
      clientId = v
    }

    if (clientId != "") {
      // Get request headers
      headers := make(map[string]string)
      for name, values := range c.Request().Header {
        // Combine multiple values for the same header into a single string
        headers[name] = values[0]
      }

      data := map[string]interface{}{
        "http_method": c.Request().Method,
        "url": c.Request().URL.String(),
        "headers": headers,
        "form_params": c.FormParams(),
        "query_params": c.queryParams(),
      }
      HistoryRepository.Record(clientId, data)
    }

    return next(c)
  }
}

