package main

import (
	// "fmt"
	"embed"
	"html/template"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
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
	engine := html.New("./templates", ".tmpl")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// fmt.Print(engine.Templates)

	app.Use(logger.New())

	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("page/index", fiber.Map{
			"Title": "testing",
		})
	})

	app.Get("/write", func(c *fiber.Ctx) error {
		return c.Render("page/write", fiber.Map{
			"Title": "testing",
		})
	})

	app.Listen(":2999")
}
