package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	util "github.com/fbold/futile.me/internal"
	"github.com/fbold/futile.me/internal/templates/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

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
	r.Handle("/static/*", http.StripPrefix("/static/", addHeaders(fs)))

	// PAGES

	r.Get("/register", util.Serve(pages.Register))
	r.Post("/register", handleRegister)
	r.Get("/login", util.Serve(pages.Login))
	r.Post("/login", handleLogin)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		// r.Use(jwtauth.Authenticator(tokenAuth))
		r.Use(authenticate())

		r.Get("/", util.Serve(pages.Home))
		r.Get("/profile", util.Serve(pages.Profile))
	})

	http.ListenAndServe(":2999", r)
}

func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		fs.ServeHTTP(w, r)
	}
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

func authenticate() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				// checks if coming from loaded page, in which case let htmx handle
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/login")
					w.WriteHeader(http.StatusSeeOther)
					return
				}

				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			if token == nil {
				// checks if coming from loaded page, in which case let htmx handle
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/login")
					w.WriteHeader(http.StatusSeeOther)
					return
				}

				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
