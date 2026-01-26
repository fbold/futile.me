package main

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fbold/futile.me/internal/templates/components"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Document struct {
	UserId  float64 `validate:"required"`
	Content string  `validate:"required"`
}

func handleDocumentCreate(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value("validator").(*validator.Validate)
	db := r.Context().Value("db").(*pgxpool.Pool)

	_, claims, _ := jwtauth.FromContext(r.Context())

	document := &Document{
		UserId:  claims["user_id"].(float64),
		Content: r.FormValue("content"),
	}

	err := v.Struct(document)
	if err != nil {
		slog.Info("user NOT valid.", "error", err)
		return
	} else {
		slog.Info("user valid", "user", document)
	}

	var createdDoc Document
	err = db.QueryRow(r.Context(), `
		INSERT INTO
			documents (user_id, content)
			VALUES ($1, $2)`, document.UserId,
		document.Content,
	).Scan(&createdDoc)

	if err != nil {
		slog.Info("error db", "error", err)
	}

	templ.Handler(components.Document(document.Content)).ServeHTTP(w, r)
}
