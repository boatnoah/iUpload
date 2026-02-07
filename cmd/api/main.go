package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/boatnoah/iupload/internal/auth"
	"github.com/boatnoah/iupload/internal/blob"
	"github.com/boatnoah/iupload/internal/db"
	"github.com/boatnoah/iupload/internal/processor"
	"github.com/boatnoah/iupload/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/joho/godotenv"
)

func main() {
	r := chi.NewRouter()

	err := godotenv.Load()
	if err != nil {
		panic(errors.New("No database creds"))
	}

	db, err := db.New(RequiredEnv("DATABASE_URL"), 30, 30, "15m")

	if err != nil {
		panic(err)
	}

	blobStore, err := blob.New(
		RequiredEnv("S3_ENDPOINT"),
		RequiredEnv("S3_REGION"),
		RequiredEnv("S3_ACCESS_KEY_ID"),
		RequiredEnv("S3_SECRET_ACCESS_KEY"),
		RequiredEnv("S3_BUCKET"),
	)
	if err != nil {
		panic(err)
	}

	store := storage.NewStorage(db)

	svc := processor.New(store, blobStore)
	auth := auth.New(store)

	app := app{
		svc,
		auth,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome"))
		})
		r.Post("/register", app.registerUserHandler)
		r.Post("/login", app.logInUser)
		r.Route("/images", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {})               // upload images
			r.Post("/{id}/transform", func(w http.ResponseWriter, r *http.Request) {}) // transform images
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Welcome to a protected route"))
			}) // transform images
		})
	})

	http.ListenAndServe(":3000", r)
}
