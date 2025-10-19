package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/fbold/futile.me/cmd/handlers"
	"github.com/fbold/futile.me/templates/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func serve(page func() templ.Component) func(w http.ResponseWriter, r *http.Request) {
	return templ.Handler(page()).ServeHTTP
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	_, dbmiddleware := connectDB()
	r.Use(dbmiddleware)

	validate := validator.New(validator.WithRequiredStructEnabled())
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "validator", validate)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	dir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(dir, "static"))
	fs := http.FileServer(filesDir)
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", serve(pages.Home))
	r.Get("/profile", serve(pages.Profile))

	r.Get("/register", serve(pages.Register))
	r.Post("/register", handlers.Register)

	http.ListenAndServe(":2999", r)
}

func connectDB() (*pgxpool.Pool, func(next http.Handler) http.Handler) {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	return dbpool, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", dbpool)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
