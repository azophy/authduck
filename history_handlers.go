package main

import (
  "log"
  "strings"
	"net/http"
  "encoding/json"
  "encoding/base64"

	"github.com/labstack/echo/v4"
)

func RegisterHistoryHandlers(app *echo.Echo) {
  app.Use(extractBasicAuthMiddleware)
  app.Use(historyRecorderMiddleware)

	app.GET("/manage/history", ServeResourceFile("resources/pages/history_detail.html"))
  app.GET("/manage/history__", historyDetailHandler)

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
                <textarea style="width:100%" rows="5" disabled>
- method: {{ .HTTPMethod }}
- url: {{ .Url }}
- query params: {{ .QueryParams }}
- headers: {{ .Headers }}
- body: {{ .Body }}
                </textarea>
                </td>
              </tr>
              {{ end }}
            {{ end }}
          </tbody>
        </table>
      </div>
  `
  res, err := RenderHTML(templ, map[string]interface{}{
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

      headerBytes, _ := json.Marshal(headers)
      formParamsBytes, _ := json.Marshal(formParams)
      queryParamsBytes, _ := json.Marshal(c.QueryParams())

      HistoryRepository.Record(
        clientId,
        c.Request().Method,
        c.Request().URL.String(),
        string(headerBytes),
        string(formParamsBytes),
        string(queryParamsBytes),
      )
    }

    return next(c)
  }
}

