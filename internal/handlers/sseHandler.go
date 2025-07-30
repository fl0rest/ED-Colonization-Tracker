package handlers

import (
	"ed-tracker/internal/logging"
	"encoding/json"
	"net/http"
	"time"
)

func SseHandler(w http.ResponseWriter, r *http.Request) {
	log := logging.Log

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming Unsupported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	_, _ = w.Write([]byte(": connected\n\n"))
	flusher.Flush()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("SSE Client Disconnected")
			return

		case <-ticker.C:
			resources, err := queries.ListResources(ctx)
			if err != nil {
				log.Errorf("SSE Query Error: %v", err)
				continue
			}

			data, err := json.Marshal(resources)
			if err != nil {
				log.Error("SSE Marshal Error:", err)
				continue
			}

			_, err = w.Write([]byte("data: " + string(data) + "\n\n"))
			if err != nil {
				log.Error("SSE Write Error:", err)
				continue
			}
			flusher.Flush()
		}
	}

}
