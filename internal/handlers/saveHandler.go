package handlers

import (
	_ "modernc.org/sqlite"
	"net/http"
	"path/filepath"
)

func SaveHandler (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
