package main

import (
  "net/http"

	"github.com/labstack/echo/v4"
)

func TokenHandler(c echo.Context) error {
    formParams, _ := c.FormParams()
		return c.JSON(http.StatusOK, map[string]interface{}{
      "params": formParams,
      "id_token": "oiiweiow",
      "access_token": "klsdjflksd",
    })
}
