package handlers

import (
	"database/sql"
	"ed-tracker/internal/db"
	"ed-tracker/internal/handlers/events"
	"ed-tracker/internal/logging"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	log := logging.Log

	var (
		dockEvent  *events.DockEvent
		depotEvent *events.DepotEvent
		unixTime   time.Time
	)

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Empty Request", http.StatusBadRequest)
		return
	}

	var rawEvents []json.RawMessage
	if err := json.Unmarshal(body, &rawEvents); err != nil {
		log.Errorf("Failed to unmarshal raw events: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, raw := range rawEvents {
		var meta events.EventMeta
		if err := json.Unmarshal(raw, &meta); err != nil {
			log.Errorf("Failed to read metadata: %v", err)
			continue
		}

		unixTime, err = time.Parse(time.RFC3339, meta.Timestamp)
		if err != nil {
			log.Errorf("Invalid timestamp: %v", err)
			continue
		}

		switch meta.Event {
		case "Docked":
			var e events.DockEvent
			if err := json.Unmarshal(raw, &e); err != nil {
				log.Errorf("Failed to parse dock event: %v", err)
				continue
			}
			dockEvent = &e

		case "ColonisationConstructionDepot":
			var e events.DepotEvent
			if err := json.Unmarshal(raw, &e); err != nil {
				log.Errorf("Failed to parse depot event: %v", err)
				continue
			}
			depotEvent = &e

		default:
			log.Infof("Skipping unknown event type: %s", meta.Event)
		}
	}

	if dockEvent == nil || depotEvent == nil {
		http.Error(w, "Invalid event bundle", http.StatusBadRequest)
		return
	}

	stationId, err := queries.GetStationId(r.Context(), sql.NullString{
		String: string(dockEvent.MarketID),
		Valid:  false,
	})

	if err != nil {
		log.Infof("Station Not Found, Adding")
		stationArguments := db.AddStationParams{
			Marketid:    int64(dockEvent.MarketID),
			Systemname:  dockEvent.StarSystem,
			Stationname: dockEvent.StationName,
		}
		if err := queries.AddStation(r.Context(), stationArguments); err != nil {
			log.Errorf("Failed to add Station: %v", err)
			http.Error(w, "Save Error", http.StatusInternalServerError)
		}
	}

	args := db.AddEventParams{
		Time:         unixTime.Unix(),
		Completion:   round2(depotEvent.Completion),
		MarketId:     int64(depotEvent.MarketID),
		StationId:    stationId,
		RawResources: string(depotEvent.Raw),
	}

	if err := queries.AddEvent(r.Context(), args); err != nil {
		log.Errorf("Failed to insert event: %w", err)
		http.Error(w, "Save Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func round2(val float64) float64 {
	return float64(int(val*100)) / 100
}
