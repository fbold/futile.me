package models

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Document struct {
	Id      int
	Created string
	Updated string
	Content string
}

func GetDocuments(r *http.Request) []Document {
	db := r.Context().Value("db").(*pgxpool.Pool)
	var documents []Document
	err := db.QueryRow(r.Context(), `
		SELECT * FROM 
			documents
		ORDER BY updated DESC
		LIMIT 4`).Scan(&documents)

	if err != nil {
		log.Println(err)
		return []Document{}
	}

	return documents
}
