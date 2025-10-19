package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type User struct {
	username string `validate:"required"`
	password string `validate:"required"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value("validator").(*validator.Validate)

	user := &User{
		username: r.FormValue("username"),
		password: r.FormValue("password"),
	}

	err := v.Struct(user)
	if err != nil {
		slog.Info("user. valid.", "error", err)
	} else {
		slog.Info("user valid")
	}
}
