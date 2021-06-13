package main

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func IncludeHTML(path string) (template.HTML, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error: could not open file: %w", err)
	}
	return template.HTML(string(data)), nil
}

func CreateTemplate() *Template {
	funcMap := template.FuncMap{
		"IncludeHTML": IncludeHTML,
	}
	t := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
	return &Template{
		templates: t,
	}
}
