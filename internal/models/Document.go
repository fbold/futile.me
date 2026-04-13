package models

import (
	"log"
	"net/http"

	"github.com/fbold/futile.me/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDocuments(r *http.Request) []sqlc.Document {
	db := r.Context().Value("db").(*pgxpool.Pool)
	q := sqlc.New(db)

	docs, err := q.GetDocuments(r.Context(), 4)
	if err != nil {
		log.Println(err)
		return []sqlc.Document{}
	}

	return docs
}
