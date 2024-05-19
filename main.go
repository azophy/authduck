package main

import (
  "io"
  "log"
  "strings"
	"net/http"
  "html/template"
  "encoding/json"
  "encoding/base64"

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

  e.Use(extractBasicAuthMiddleware)
  e.Use(historyRecorderMiddleware)

	e.GET("/", ServeResourceFile("resources/pages/index.html"))
	e.GET("/assets/*", ServeResourceFolder("resources"))

	e.GET("/.well-known/certs", func (c echo.Context) error {
		return c.JSON(http.StatusOK, PublicJWKS)
  })

  e.Renderer = NewTemplateRenderer()
	e.GET("/manage/history", ServeResourceFile("resources/pages/history_detail.html"))
  e.GET("/manage/history__", historyDetailHandler)

  RegisterGeneralOAuthModule(e)

	e.Logger.Fatal(e.Start(":" + APP_PORT))
}

type Template struct {
  templates *template.Template
}

// Render method to render the template string
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  tmpl, err := template.New("htmlTemplate").Parse(name)
  if err != nil {
    return err
  }
  return tmpl.Execute(w, data)
}

func NewTemplateRenderer() *Template {
  renderer := &Template{
    templates: template.Must(template.New("T").Parse("")),
  }

  return renderer
}

func renderHTML(htmlTemplate string, data interface{}) (string, error) {
  // Parse and execute the template
  tmpl, err := template.New("htmlTemplate").Parse(htmlTemplate)
  if err != nil {
    return "", err
  }

  // Create a buffer to hold the executed template
  var renderedContent strings.Builder
  if err := tmpl.Execute(&renderedContent, data); err != nil {
    return "", err
  }

  // Return the rendered HTML
  return renderedContent.String(), nil
}

func historyDetailHandler(c echo.Context) error {
  clientId := c.QueryParam("id")
  from := c.QueryParam("from")
  histories, err := HistoryRepository.All(clientId, from)
  if err != nil {
    return err
  }

  templ := `
      <div id="history-list">
        <form hx-get="/manage/history__" hx-target="#history-list">
          <label for="">client_id</label>
          <input type="text" name="id" value="{{ .client_id }}">
          <input type="hidden" name="from" value="0">
          <button type="submit">get</button>
        </form>

        <table>
          <thead>
            <tr>
              <th>timestamp</th>
              <th>data</th>
            </tr>
          </thead>
          <tbody>
          {{ if len .histories | ge 0 }}
              empty data for client "{{ .client_id }}"
            {{ else }}
              {{ range .histories }}
              <tr>
                <td>{{ .Timestamp }}</td>
                <td>
                <textarea style="width:100%" rows="5" disabled>{{ .Data }}</textarea>
                </td>
              </tr>
              {{ end }}
            {{ end }}
          </tbody>
        </table>
      </div>
  `
  res, err := renderHTML(templ, map[string]interface{}{
    "histories": histories,
    "client_id": clientId,
  })
  if err != nil {
    log.Printf("found err %v\n", err)
    return err
  }

  return c.HTML(http.StatusOK, res)
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
      log.Printf("got client_id=%v\n", clientId)

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

