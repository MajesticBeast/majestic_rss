package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/majesticbeast/majestic_rss/internal/database"

	"github.com/go-chi/chi"
)

func (s *apiConfig) v1Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/readiness", handleReadiness)
	r.Get("/err", handleErr)
	r.Get("/feeds", s.handleGetFeeds)
	r.Get("/admin", s.handleAdmin)
	r.Get("/admin/feeds/update/{id}", s.handleUpdateFeedForm)
	r.Get("/admin/feeds/{id}", s.handleGetFeed)

	r.Post("/feeds", s.handleCreateFeed)
	r.Delete("/admin/feeds/delete/{id}", s.handleDeleteFeed)
	r.Put("/admin/feeds/update/{id}", s.handleUpdateFeed)
	return r
}

func (s *apiConfig) handleGetFeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid ID sent: "+err.Error())
		return
	}

	// Get feed from database
	feed, err := s.DB.GetFeedByID(r.Context(), int32(idInt))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to find ID: "+err.Error())
		return
	}

	// Generate single row for a feed
	tmpl := template.Must(template.ParseFiles(
		"templates/feed.html",
	))

	err = tmpl.Execute(w, feed)
	if err != nil {
		fmt.Println("Unable to execute template: ", err)
		return
	}
}

func (s *apiConfig) handleUpdateFeedForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Received invalid id: ", id)
		respondWithError(w, http.StatusInternalServerError, "Unable to load edit form: "+err.Error())
		return
	}

	// Get feed from database
	feed, err := s.DB.GetFeedByID(r.Context(), int32(idInt))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to load edit form: "+err.Error())
		return
	}

	// Generate edit form template and return it
	tmpl := template.Must(template.ParseFiles(
		"templates/edit.html",
	))

	err = tmpl.Execute(w, feed)
	if err != nil {
		fmt.Println("Unable to execute template: ", err)
		return

	}

}

func (s *apiConfig) handleUpdateFeed(w http.ResponseWriter, r *http.Request) {
	// parse form data
	type parameters struct {
		ID         string
		Name       string
		FeedUrl    string
		WebhookUrl string
	}

	params := parameters{}
	params.ID = r.PostFormValue("ID")
	params.Name = r.PostFormValue("Name")
	params.FeedUrl = r.PostFormValue("FeedUrl")
	params.WebhookUrl = r.PostFormValue("WebhookUrl")

	idInt, err := strconv.Atoi(params.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update feed: "+err.Error())
		return
	}

	err = s.DB.UpdateFeed(r.Context(), database.UpdateFeedParams{
		ID:         int32(idInt),
		Name:       params.Name,
		FeedUrl:    params.FeedUrl,
		WebhookUrl: params.WebhookUrl,
		UpdatedAt:  time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update feed: "+err.Error())
		return
	}

	feed, err := s.DB.GetFeedByID(r.Context(), int32(idInt))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update feed: "+err.Error())
		return
	}

	// Generate single row for a feed
	tmpl := template.Must(template.ParseFiles(
		"templates/feed.html",
	))

	err = tmpl.Execute(w, feed)
	if err != nil {
		fmt.Println("Unable to execute template: ", err)
		return
	}
}

func (s *apiConfig) handleDeleteFeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Received invalid id: ", id)
		respondWithError(w, http.StatusInternalServerError, "Unable to delete feed: "+err.Error())
		return
	}

	err = s.DB.DeleteFeed(r.Context(), int32(idInt))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete feed: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *apiConfig) handleAdmin(w http.ResponseWriter, r *http.Request) {
	// display template
	tmpl := template.Must(template.ParseFiles(
		"templates/admin.html",
		"templates/header.html",
	))

	// get feeds
	feeds, err := s.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to list feeds: "+err.Error())
		return
	}

	// pass feeds to template
	type templateData struct {
		Feeds []database.Feed
	}

	data := templateData{
		Feeds: feeds,
	}

	err = tmpl.Execute(w, data)
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
		Name       string    `json:"name"`
		FeedURL    string    `json:"feed_url"`
		WebhookURL string    `json:"webhook_url"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}
	params := parameters{}

	params.Name = r.PostFormValue("name")
	params.FeedURL = r.PostFormValue("feed_url")
	params.WebhookURL = r.PostFormValue("webhook_url")
	params.CreatedAt = time.Now().UTC()
	params.UpdatedAt = time.Now().UTC()

	err := s.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		Name:       params.Name,
		FeedUrl:    params.FeedURL,
		WebhookUrl: params.WebhookURL,
		CreatedAt:  params.CreatedAt,
		UpdatedAt:  params.UpdatedAt,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create feed: "+err.Error())
		return
	}

	feed, err := s.DB.GetFeedByName(r.Context(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update feed: "+err.Error())
		return
	}

	// Generate single row for a feed
	tmpl := template.Must(template.ParseFiles(
		"templates/feed.html",
	))

	err = tmpl.Execute(w, feed)
	if err != nil {
		fmt.Println("Unable to execute template: ", err)
		return
	}
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
