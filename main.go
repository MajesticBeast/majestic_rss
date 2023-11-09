package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/majesticbeast/majestic_rss/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")
	dbconn := os.Getenv("DBCONN")
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		log.Fatal("Unable to set concurrency: ", err)
	}
	interval, err := strconv.Atoi(os.Getenv("INTERVAL"))
	if err != nil {
		log.Fatal("Unable to set interval: ", err)
	}

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

	collectionConcurrency := concurrency
	collectionInterval := time.Duration(interval) * time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Println("Listening at: " + port)
	log.Fatal(server.ListenAndServe())
}
