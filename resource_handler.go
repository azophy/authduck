package main

import (
  "io"
  "io/fs"
  "log"
	"embed"
  "errors"
  "regexp"
  "strings"
  "net/http"
  "html/template"

	"github.com/labstack/echo/v4"
)

//go:embed resources
var embededFiles embed.FS
var templateModifierRegex = regexp.MustCompile(`#([a-z_0-9]+)`)

// Define the template registry struct
type TemplateRegistry struct {
  templates map[string]*template.Template
  baseTemplatePaths map[string]string
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, templateName string, data interface{}, c echo.Context) error {
  name := templateModifierRegex.ReplaceAllString(templateName, "")
  log.Printf("name: %v\n", name)

  tmpl, ok := t.templates[name]
  if !ok {
    err := errors.New("Template not found -> " + name)
    log.Println(err)
    return err
  }

  finalData := map[string]interface{}{
    "config": Config,
  }

  switch dataMap := data.(type) {
    case map[string]interface{}:
      for k,v := range(dataMap) {
        finalData[k] = v
      }
    case echo.Map :
      for k,v := range(dataMap) {
        finalData[k] = v
      }
    default:
      return errors.New("supplied data is not a map-like type")
  }

  renderName := t.baseTemplatePaths[name]
  res := templateModifierRegex.FindStringSubmatch(templateName)
  if len(res) > 1 && res[1] == "partial" {
    renderName = "body"
  }
  log.Printf("rendername: %v\n", renderName)
  err := tmpl.ExecuteTemplate(w, renderName, finalData)
  if err != nil {
    log.Println(err)
    return err
  }

  return nil
}

func SetupTemplateRegistry(parentPath string) *TemplateRegistry {
  files, err := fs.Glob(embededFiles, parentPath)
  if err != nil {
      log.Println(err)
  }

  log.Printf("found %v files\n", len(files))
  templates := make(map[string]*template.Template)
  baseTemplatePaths := make(map[string]string)
  for _, filePath := range files {
    log.Printf("processing template file %v\n", filePath)
    tmpl, err := template.ParseFS(embededFiles, filePath)
    if err != nil {
        log.Println(err)
        continue
    }

    // by default render partial template
    baseTemplatePath := "body"
    // get basepath definition
    var renderedContent strings.Builder
    err = tmpl.ExecuteTemplate(&renderedContent, "base_template_path", nil)
    if err != nil {
        log.Println(err)
        log.Println("skipping...")
    } else {
        templatePath := renderedContent.String()
        res,_ := template.ParseFS(embededFiles, filePath, templatePath)
        if err != nil {
            log.Println(err)
            continue
        }

        tmpl = res
        baseTemplatePath = templatePath
    }

    templates[filePath] = tmpl
    baseTemplatePaths[filePath] = baseTemplatePath
  }

  return &TemplateRegistry{
    templates: templates,
    baseTemplatePaths: baseTemplatePaths,
  }
}

func ServeResourceTemplate(path string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, path, nil)
  }
}

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
