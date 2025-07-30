package handlers

import (
	"ed-tracker/internal/db"
	"ed-tracker/internal/logging"
	"net/http"
	"text/template"

	_ "modernc.org/sqlite"
)

var (
	indexTmpl *template.Template
	queries   *db.Queries
)

func Init(q *db.Queries) {
	queries = q
	var err error
	indexTmpl, err = template.ParseFiles("internal/templates/index.tmpl")
	if err != nil {
		logging.Log.Errorf("Failed to parse template: %v", err)
	}
}

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
