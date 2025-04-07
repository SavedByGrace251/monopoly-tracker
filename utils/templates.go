package utils

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w http.ResponseWriter, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}
