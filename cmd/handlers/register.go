package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	username string `validate:"required"`
	password string `validate:"required"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value("validator").(*validator.Validate)
	db := r.Context().Value("db").(*pgxpool.Pool)

	user := &User{
		username: r.FormValue("username"),
		password: r.FormValue("password"),
	}

	db.QueryRow(context.Background(), "INSERT INTO users VALUES ()")

	err := v.Struct(user)
	if err != nil {
		slog.Info("user. valid.", "error", err)
	} else {
		slog.Info("user valid")
	}
}
