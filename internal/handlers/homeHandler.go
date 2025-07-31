package handlers

import (
	"ed-tracker/internal/logging"
	"net/http"

	_ "modernc.org/sqlite"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.Log

	resources, err := queries.ListEvents(ctx)
	if err != nil {
		log.Errorf("Error listing Resources: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := indexTmpl.Execute(w, resources); err != nil {
		log.Errorf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
