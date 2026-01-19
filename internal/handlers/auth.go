package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	util "github.com/fbold/futile.me/internal"
	"github.com/fbold/futile.me/internal/templates/pages"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/sqids/sqids-go"
)

type User struct {
	Username        string  `validate:"required"`
	Email           *string `validate:"omitempty,email"`
	Password        string  `validate:"required"`
	PasswordConfirm string  `validate:"required,eqfield=Password"`
}

func auth() {

	return
}

func Register(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value("validator").(*validator.Validate)
	db := r.Context().Value("db").(*pgxpool.Pool)

	user := &User{
		Username:        r.FormValue("username"),
		Email:           util.NullString(r.FormValue("email")),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("passwordconfirm"),
	}

	err := v.Struct(user)
	if err != nil {
		slog.Info("user NOT valid.", "error", err)
		return
	} else {
		slog.Info("user valid", "user", user)
	}

	var createdUser User
	err = db.QueryRow(context.Background(), `
		INSERT INTO
			users (username, email, password)
			VALUES ($1, $2, $3)`,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&createdUser)

	if err != nil {
		slog.Info("error db", "error", err)
	}
}

type UserLogin struct {
	Username string
	Password string
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*pgxpool.Pool)

	userRequest := UserLogin{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	var username string
	var password string

	var err = db.QueryRow(context.Background(), `
		SELECT username, password FROM users WHERE username = $1`,
		userRequest.Username,
	).Scan(&username, &password)

	if err != nil {
		slog.Error("No user matching", err)
		http.Error(w, "No user", http.StatusUnauthorized)
		return
	}

	if userRequest.Password != password {
		slog.Error("Wrong password", err)
		http.Error(w, "Wrong", http.StatusUnauthorized)
		return
	}

	slog.Info("user found", "user", userRequest)

	w.Header().Add("HX-Location", "/")
	util.Serve(pages.Home)
}
