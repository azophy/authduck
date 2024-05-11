package main

import (
  "net/http"

	"github.com/labstack/echo/v4"
)

func TokenHandler(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
      "params": c.Request().URL.Query(),
    })
}
