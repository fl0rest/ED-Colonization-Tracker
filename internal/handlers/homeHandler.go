package handlers

import (
	"ed-tracker/internal/db"
	"ed-tracker/internal/logging"
	"net/http"
	"path/filepath"
	"strings"
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

	data := IndexPageData{
		LastUpdate: time.Unix(event.Time, 0).Format("2006-01-02 15:04:05"),
		Resources:  resources,
		Progress:   float64(event.Completion * 100),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := indexTmpl.Execute(w, data); err != nil {
		log.Errorf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	assetPath := strings.TrimPrefix(r.URL.Path, "/static/")

	if strings.Contains(assetPath, "..") {
		http.NotFound(w, r)
		return
	}

	if !strings.Contains(assetPath, ".svg") {
		http.NotFound(w, r)
		return
	}

	fullPath := filepath.Join("static", assetPath)

	http.ServeFile(w, r, fullPath)
}
