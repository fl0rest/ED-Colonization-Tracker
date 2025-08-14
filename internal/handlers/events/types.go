package events

import "encoding/json"

type EventMeta struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
}

type DepotEvent struct {
	Timestamp  string          `json:"timestamp"`
	Event      string          `json:"event"`
	MarketID   int             `json:"MarketID"`
	Completion float64         `json:"ConstructionProgress"`
	Raw        json.RawMessage `json:"ResourcesRequired"`
}

type DockEvent struct {
	Timestamp   string `json:"timestamp"`
	Event       string `json:"event"`
	MarketID    int    `json:"MarketID"`
	StarSystem  string `json:"StarSystem"`
	StationName string `json:"StationName"`
}
