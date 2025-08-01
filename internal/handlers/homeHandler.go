package handlers

import (
	"ed-tracker/internal/db"
	"ed-tracker/internal/logging"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type IndexPageData struct {
	LastUpdate string
	Progress   float64
	Resources  []db.Resource
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.Log

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	resources, err := queries.ListResources(ctx)
	if err != nil {
		log.Errorf("Error listing Resources: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	event, err := queries.GetLatestEvent(ctx)
	if err != nil {
		log.Errorf("Error fetching event: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	var lastUpdate string
	if len(resources) > 0 {
		lastUpdate = time.Unix(resources[0].Time, 0).Format("2006-01-02 15:04:05")
	}

	data := IndexPageData{
		LastUpdate: lastUpdate,
		Resources:  resources,
		Progress:   event.Completion,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := indexTmpl.Execute(w, data); err != nil {
		log.Errorf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
