package main

import (
	// "fmt"
	"embed"
	"html/template"
	"io"

	"github.com/a-h/templ"
	"github.com/fbold/futile.me/templates/pages"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed static/*
var static embed.FS

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.templ.html")),
	}
}

func main() {
	app := fiber.New(fiber.Config{})

	app.Use(logger.New())

	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return adaptor.HTTPHandler(templ.Handler(pages.Home()))(c)
	})

	app.Get("/write", func(c *fiber.Ctx) error {
		return c.Render("page/write", fiber.Map{
			"Title": "testing",
		})
	})

	app.Listen(":2999")
}
