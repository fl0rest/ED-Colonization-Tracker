package handlers

import (
	_ "modernc.org/sqlite"
	"net/http"
	"path/filepath"
)

func FetchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
