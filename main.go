package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/majesticbeast/majestic_rss/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (s *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apikey := r.Header.Get("Authorization")
		apikey = strings.TrimPrefix(apikey, "ApiKey ")

		user, err := s.DB.GetUserByAPIKey(r.Context(), apikey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unauthorized")
			return
		}

		handler(w, r, user)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	dbconn := os.Getenv("DBCONN")

	db, err := sql.Open("postgres", dbconn)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	apiConfig := &apiConfig{
		DB: dbQueries,
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Mount("/v1", apiConfig.v1Router())

	server := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Println("Listening at: " + port)
	log.Fatal(server.ListenAndServe())
}
