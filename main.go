package main

import (
  "log"
	"embed"
	"io/fs"
  "strings"
	"net/http"
  "encoding/json"
  "encoding/base64"

	"github.com/labstack/echo/v4"
)

//go:embed resources
var embededFiles embed.FS

func main() {
  APP_PORT := GetEnvOrDefault("APP_PORT", "3000")
  err := InitiateGlobalVars()
  if err != nil {
    log.Printf("failed initiating global vars: %s\n", err)
    return
  }

	e := echo.New()

  e.Use(extractBasicAuthMiddleware)
  e.Use(historyRecorderMiddleware)

  assetHandler := http.FileServer(getFileSystem())
	e.GET("/", serveEmbededFile("resources/pages/index.html"))
	e.GET("/assets/*", echo.WrapHandler(assetHandler))

	e.GET("/.well-known/certs", func (c echo.Context) error {
		return c.JSON(http.StatusOK, PublicJWKS)
  })

  e.GET("/manage/history/:client_id/:from", func (c echo.Context) error {
    clientId := c.Param("client_id")
    from := c.Param("from")
    histories, err := HistoryRepository.All(clientId, from)
    if err != nil {
      return err
    }
    return c.JSON(http.StatusOK, histories)
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

func extractBasicAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
    // Get the Authorization header
    authHeader := c.Request().Header.Get("Authorization")
    if authHeader == "" {
      return next(c)
    }

    // Check if the Authorization header starts with "Basic "
    if !strings.HasPrefix(authHeader, "Basic ") {
      return next(c)
    }

    // Decode the base64 part of the header
    base64Credentials := strings.TrimPrefix(authHeader, "Basic ")
    credentials, err := base64.StdEncoding.DecodeString(base64Credentials)
    if err != nil {
      return next(c)
    }

    // Split the credentials into username and password
    creds := strings.SplitN(string(credentials), ":", 2)
    if len(creds) != 2 {
      return next(c)
    }

    c.Set("basic-auth-username", creds[0])
    c.Set("basic-auth-password", creds[1])

    return next(c)
  }
}



func historyRecorderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
    var clientId string

    if v := c.Get("basic-auth-username"); v != nil {
      clientId = v.(string)
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

      formParams, _ := c.FormParams()

      data := map[string]interface{}{
        "http_method": c.Request().Method,
        "url": c.Request().URL.String(),
        "headers": headers,
        "form_params": formParams,
        "query_params": c.QueryParams(),
      }
      bytes, _ := json.Marshal(data)
      HistoryRepository.Record(clientId, string(bytes))
    }

    return next(c)
  }
}

