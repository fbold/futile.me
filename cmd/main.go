package main

import (
	"github.com/a-h/templ"
	"github.com/fbold/futile.me/cmd/handlers"
	"github.com/fbold/futile.me/templates/pages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func usingFiber(page func() templ.Component) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return adaptor.HTTPHandler(templ.Handler(page()))(c)
	}
}

func main() {
	app := fiber.New(fiber.Config{})

	app.Use(logger.New())

	app.Static("/static", "./static")

	validate := validator.New(validator.WithRequiredStructEnabled())

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("validator", validate)
		return c.Next()
	})

	app.Get("/", usingFiber(pages.Home))
	app.Get("/profile", usingFiber(pages.Profile))

	app.Get("/register", usingFiber(pages.Register))
	app.Post("/register", handlers.Register)

	app.Listen(":2999")
}

func connectDB() {

}
