package main

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fbold/futile.me/internal/sqlc"
	"github.com/fbold/futile.me/internal/templates/components"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateDocument struct {
	UserID  int32  `validate:"required"`
	Content string `validate:"required"`
}

func handleDocumentCreate(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value("validator").(*validator.Validate)
	db := r.Context().Value("db").(*pgxpool.Pool)
	q := sqlc.New(db)

	_, claims, _ := jwtauth.FromContext(r.Context())

	doc := &CreateDocument{
		UserID:  int32(claims["user_id"].(float64)),
		Content: r.FormValue("content"),
	}

	err := v.Struct(doc)
	if err != nil {
		slog.Info("document NOT valid.", "error", err)
		return
	}
	slog.Info("document valid", "user", doc)

	_, err = q.CreateDocument(r.Context(), sqlc.CreateDocumentParams{
		UserID:  doc.UserID,
		Content: doc.Content,
		Private: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		slog.Info("error db", "error", err)
	}

	templ.Handler(components.Document(doc.Content)).ServeHTTP(w, r)
}
