package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sqids/sqids-go"
)

type User struct {
	username        string `validate:"required"`
	email           string `validate:"email"`
	password        string `validate:"required"`
	passwordconfirm string `validate:"required,eqfield=password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value("validator").(*validator.Validate)
	db := r.Context().Value("db").(*pgxpool.Pool)

	user := &User{
		username: r.FormValue("username"),
		email:    r.FormValue("email"),
		password: r.FormValue("password"),
	}

	err := v.Struct(user)
	if err != nil {
		slog.Info("user. valid.", "error", err)
	} else {
		slog.Info("user valid")
	}

	var createdUser User
	err = db.QueryRow(context.Background(), `
		INSERT INTO users (username, email, password) VALUES (
			$1, $2, $3
		)`, user.email, user.email, user.password).Scan(&createdUser)

	if err != nil {
		slog.Info("error db", "error", err)
	}

	slog.Info("user valid", "user", createdUser)
}
