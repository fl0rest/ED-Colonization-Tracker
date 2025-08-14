package handlers

import (
	"database/sql"
	"ed-tracker/internal/db"
	"ed-tracker/internal/logging"
	"encoding/json"
	"net/http"
	"time"
)

type dataPayload struct {
	Resources   []db.Resource `json:"resources"`
	Progress    float64       `json:"progress"`
	LastUpdated string        `json:"lastUpdated"`
}

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
			resources, err := queries.ListResources(ctx, sql.NullString{
				String: string(0),
				Valid:  false,
			})
			if err != nil {
				log.Errorf("SSE Query Error: %v", err)
				continue
			}

			event, err := queries.GetLatestEvent(ctx)
			if err != nil {
				log.Error("SSE Query Error", err)
				continue
			}
			lastUpdate := time.Unix(event.Time, 0).Format("2006-01-02 15:04:05")

			payload := dataPayload{
				Resources:   resources,
				Progress:    float64(event.Completion * 100),
				LastUpdated: lastUpdate,
			}

			data, err := json.Marshal(payload)
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
