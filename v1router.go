package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/majesticbeast/majestic_rss/internal/database"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func (s *apiConfig) v1Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/readiness", handleReadiness)
	r.Get("/err", handleErr)
	r.Get("/users", s.middlewareAuth(s.handleGetUser))
	r.Get("/feeds", s.handleGetFeeds)
	r.Get("/feed_follows", s.middlewareAuth(s.handleGetFeedFollows))

	r.Post("/users", s.handleCreateUser)
	r.Post("/feeds", s.middlewareAuth(s.handleCreateFeed))
	r.Post("/feed_follows", s.middlewareAuth(s.handleCreateFeedFollow))

	r.Delete("/feed_follows/{feedFollowID}", s.handleDeleteFeedFollow)

	return r
}

func (s *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := s.DB.GetFeedFollowsByUserID(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to list feed follows: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, feedFollows)
}

func (s *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowID, err := uuid.Parse(chi.URLParam(r, "feedFollowID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid FeedFollowID")
		return
	}

	err = s.DB.DeleteFeedFollow(r.Context(), feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete feed follow: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID string `json:"feed_id"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body).Decode(&params)
	if err := decoder; err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	feedID, err := uuid.Parse(params.FeedID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid FeedID")
		return
	}

	feedFollow, err := s.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feedID,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create feed follow: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, feedFollow)
}

func (s *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := s.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to list feeds: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}

func (s *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
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

	feed, err := s.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	feedFollow, err := s.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create feed: "+err.Error())
		return
	}

	compoundStruct := struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{
		feed,
		feedFollow,
	}

	respondWithJSON(w, http.StatusOK, compoundStruct)
}

func (s *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, user)
}

func (s *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body).Decode(&params)
	if err := decoder; err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user, err := s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
