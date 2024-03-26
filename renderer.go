package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Renderer struct {
	templates *template.Template
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

func NewRenderer() *Renderer {
	return &Renderer{
		templates: template.Must(template.ParseGlob("views/*.gohtml")),
	}
}
