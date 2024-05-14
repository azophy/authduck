package main

import (
	"embed"
	"io/fs"
  "net/http"

	"github.com/labstack/echo/v4"
)

//go:embed resources
var embededFiles embed.FS

func ServeResourceFile(path string) echo.HandlerFunc {
  return echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.ServeFileFS(w, r, embededFiles, path)
  }))
}

func ServeResourceFolder(path string) echo.HandlerFunc {
	fsys, err := fs.Sub(embededFiles, path)
	if err != nil {
		panic(err)
	}

  resourceFS := http.FS(fsys)
  resourceHandler := http.FileServer(resourceFS)
  return echo.WrapHandler(resourceHandler)
}
