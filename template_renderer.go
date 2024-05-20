package main

import (
  "io"
  "strings"
  "html/template"

	"github.com/labstack/echo/v4"
)


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

func RenderHTML(htmlTemplate string, data interface{}) (string, error) {
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

