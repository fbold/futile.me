package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/fbold/futile.me/internal/sqlc"
	"github.com/fbold/futile.me/internal/templates/pages"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func handleServe(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*pgxpool.Pool)
	q := sqlc.New(db)

	docs, err := q.GetDocuments(r.Context(), 4)
	if err != nil {
		docs = []sqlc.Document{}
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	isLoggedIn := claims != nil && claims["user_id"] != nil

	templ.Handler(pages.Home(docs, isLoggedIn)).ServeHTTP(w, r)
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	username := ""
	if claims != nil {
		if u, ok := claims["username"].(string); ok {
			username = u
		}
	}

	templ.Handler(pages.Profile(username, true)).ServeHTTP(w, r)
}
