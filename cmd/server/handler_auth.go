package main

import (
	"log/slog"
	"net/http"
	"time"

	util "github.com/fbold/futile.me/internal"
	"github.com/fbold/futile.me/internal/sqlc"
	"github.com/fbold/futile.me/internal/templates/pages"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RegisterUser struct {
	Username        string  `validate:"required"`
	Email           *string `validate:"omitempty,email"`
	Password        string  `validate:"required"`
	PasswordConfirm string  `validate:"required,eqfield=Password"`
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value("validator").(*validator.Validate)
	db := r.Context().Value("db").(*pgxpool.Pool)
	q := sqlc.New(db)

	user := &RegisterUser{
		Username:        r.FormValue("username"),
		Email:           util.NullString(r.FormValue("email")),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("passwordconfirm"),
	}

	err := v.Struct(user)
	if err != nil {
		slog.Info("user NOT valid.", "error", err)
		return
	}
	slog.Info("user valid", "user", user)

	var email pgtype.Text
	if user.Email != nil {
		email.String = *user.Email
		email.Valid = true
	}

	_, err = q.CreateUser(r.Context(), sqlc.CreateUserParams{
		Username: user.Username,
		Email:    email,
		Password: user.Password,
	})
	if err != nil {
		slog.Info("error db", "error", err)
	}

	w.Header().Add("HX-Location", "/")
	util.Serve(pages.Login)
}

type LoginUser struct {
	Username string
	Password string
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*pgxpool.Pool)
	q := sqlc.New(db)

	_, claims, _ := jwtauth.FromContext(r.Context())

	slog.Info("jwt", "jwt", claims)

	userReq := LoginUser{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	dbUser, err := q.GetUserByUsername(r.Context(), userReq.Username)
	if err != nil {
		slog.Error("No user matching", err)
		http.Error(w, "No user", http.StatusUnauthorized)
		return
	}

	if userReq.Password != dbUser.Password {
		slog.Error("Wrong password", err)
		http.Error(w, "Wrong", http.StatusUnauthorized)
		return
	}

	_, jwtString, _ := tokenAuth.Encode(map[string]interface{}{
		"username": dbUser.Username,
		"user_id":  dbUser.ID,
		"exp":      time.Now().Add(30 * time.Minute).Unix(),
	})

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    jwtString,
		Quoted:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Location", "/")

	docs, _ := q.GetDocuments(r.Context(), 4)
	pages.Home(docs)
}
