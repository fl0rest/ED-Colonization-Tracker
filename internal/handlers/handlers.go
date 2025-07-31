package handlers

import (
	"ed-tracker/internal/db"
	"ed-tracker/internal/logging"
	"text/template"
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
