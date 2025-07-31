package handlers

import (
	"ed-tracker/internal/db"
	"ed-tracker/internal/logging"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"time"
)

type Event struct {
	Timestamp  string          `json:"timestamp"`
	Completion float64         `json:"ConstructionProgress"`
	Raw        json.RawMessage `json:"ResourcesRequired"`
	MarketId   int             `json:"MarketID"`
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	log := logging.Log

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		http.Error(w, "Empty Request", http.StatusBadRequest)
		return
	}

	var jon Event
	if err := json.Unmarshal(body, &jon); err != nil {
		log.Errorf("JSON Error: %v", err)
		return
	}

	unixTime, err := time.Parse(time.RFC3339, jon.Timestamp)
	if err != nil {
		log.Errorf("Error parsing timestamp: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	saveArgs := db.AddEventParams{
		RawText:    string(jon.Raw),
		Completion: math.Round(jon.Completion*100) / 100,
		MarketId:   int64(jon.MarketId),
		Time:       unixTime.Unix(),
	}
	_, err = queries.AddEvent(r.Context(), saveArgs)
}
