package main

import (
	"encoding/json"
	"github.com/majesticbeast/majestic_rss/internal/database"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func (s *apiConfig) v1Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/readiness", handleReadiness)
	r.Get("/err", handleErr)
	r.Get("/feeds", s.handleGetFeeds)

	r.Post("/feeds", s.handleCreateFeed)

	return r
}

func (s *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := s.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to list feeds: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}

func (s *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body).Decode(&params)
	if err := decoder; err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err := s.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		Name:      params.Name,
		Url:       params.URL,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create feed: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, "feed added")
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
