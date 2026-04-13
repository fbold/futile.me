package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/fbold/futile.me/internal/sqlc"
	"github.com/fbold/futile.me/internal/templates/pages"
	"github.com/jackc/pgx/v5/pgxpool"
)

func handleServe(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*pgxpool.Pool)
	q := sqlc.New(db)

	docs, err := q.GetDocuments(r.Context(), 4)
	if err != nil {
		docs = []sqlc.Document{}
	}

	templ.Handler(pages.Home(docs)).ServeHTTP(w, r)
}
