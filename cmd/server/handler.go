package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/fbold/futile.me/internal/models"
	"github.com/fbold/futile.me/internal/templates/pages"
)

func handleServe(w http.ResponseWriter, r *http.Request) {

	var documents = models.GetDocuments(r)

	templ.Handler(pages.Home(documents)).ServeHTTP(w, r)
}
