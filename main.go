package main

import (
	"embed"

	"github.com/a-h/templ"
	"github.com/fbold/futile.me/templates/pages"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed static/*
var static embed.FS

func usingFiber(page func() templ.Component) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return adaptor.HTTPHandler(templ.Handler(page()))(c)
	}
}

func main() {
	app := fiber.New(fiber.Config{})

	app.Use(logger.New())

	app.Static("/static", "./static")

	app.Get("/", usingFiber(pages.Home))
	app.Get("/write", usingFiber(pages.Write))

	app.Listen(":2999")
}
