package handlers

import (
	"context"
	"ed-tracker/internal/db"
	"ed-tracker/internal/logging"
	"encoding/json"
	"net/http"

	_ "modernc.org/sqlite"
)

type Response struct {
	Message string `json:"message"`
}

func FetchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	log := logging.Log

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	factory := db.DBFactory{}

	queries, rawDB, err := factory.Connect(ctx)
	if err != nil {
		log.Errorf("db connection failed: %v", err)
	}
	defer rawDB.Close()

	resources, err := queries.ListResources(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resources); err != nil {
		log.Errorf("Issue Encoding Response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
