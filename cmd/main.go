package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	util "github.com/fbold/futile.me/internal"
	"github.com/fbold/futile.me/internal/handlers"
	"github.com/fbold/futile.me/internal/templates/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// DATABASE
	dbpool, dbmiddleware := connectDB()
	defer dbpool.Close()
	r.Use(dbmiddleware)

	// VALIDATION
	validate := validator.New(validator.WithRequiredStructEnabled())
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "validator", validate)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// STATIC FILES
	dir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(dir, "static"))
	fs := http.FileServer(filesDir)
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// PAGES
	r.Get("/", util.Serve(pages.Home))
	r.Get("/profile", util.Serve(pages.Profile))

	r.Get("/register", util.Serve(pages.Register))
	r.Post("/register", handlers.Register)
	r.Get("/login", util.Serve(pages.Login))
	r.Post("/login", handlers.Login)

	http.ListenAndServe(":2999", r)
}

func connectDB() (*pgxpool.Pool, func(next http.Handler) http.Handler) {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	return dbpool, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", dbpool)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
