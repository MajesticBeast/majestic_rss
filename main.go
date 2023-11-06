package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/majesticbeast/majestic_rss/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	dbconn := os.Getenv("DBCONN")

	db, err := sql.Open("sqlite3", dbconn)
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

	const collectionConcurrency = 10
	const collectionInterval = 5 * time.Second
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Println("Listening at: " + port)
	log.Fatal(server.ListenAndServe())
}
