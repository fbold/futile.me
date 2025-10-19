package handlers

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	username string `validate:"required"`
	password string `validate:"required"`
}

func Register(c *fiber.Ctx) error {
	v := c.Locals("validator").(*validator.Validate)

	user := &User{
		username: c.FormValue("username"),
		password: c.FormValue("password"),
	}

	err := v.Struct(user)
	if err != nil {
		slog.Info("user. valid.", "error", err)
	} else {
		slog.Info("user valid")
	}

	return nil
}
