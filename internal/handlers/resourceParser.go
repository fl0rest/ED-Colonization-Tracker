package handlers

import (
	"context"
	"database/sql"
	"ed-tracker/internal/db"
	"ed-tracker/internal/logging"
	"encoding/json"
	"fmt"
)

type Resource struct {
	NameLocalised  string `json:"Name_Localised"`
	RequiredAmount int    `json:"RequiredAmount"`
	ProvidedAmount int    `json:"ProvidedAmount"`
	Payment        int    `json:"Payment"`
}

type ParsedEvent struct {
	Timestamp int        `json:"time"`
	Resources []Resource `json:"ResourcesRequired"`
}

func ParseLatestEvent(ctx context.Context, q *db.Queries) error {
	log := logging.Log

	eventRaw, err := q.GetLatestEvent(ctx)
	if err != nil {
		log.Errorf("Error getting event: %v", err)
		return err
	}

	var event ParsedEvent
	if err := json.Unmarshal([]byte(eventRaw.RawText), &event); err != nil {
		log.Errorf("JSON Error: %v", err)
		return err
	}

	for _, resource := range event.Resources {
		id, err := GetResourceId(ctx, resource.NameLocalised)
		if err != nil || id == 0 {
			log.Warnf("Resource ID not found for %s: %v, skipping", resource.NameLocalised, err)
			continue
		}

		diff := int64(resource.RequiredAmount) - int64(resource.ProvidedAmount)

		err = q.UpsertResource(ctx, db.UpsertResourceParams{
			ID:       id,
			Eventid:  eventRaw.ID,
			Name:     resource.NameLocalised,
			Required: int64(resource.RequiredAmount),
			Provided: int64(resource.ProvidedAmount),
			Payment:  int64(resource.Payment),
			Diff:     diff,
			Time:     eventRaw.Time,
		})

		if err != nil {
			log.Errorf("Error inserting: %v", err)
			return err
		}
	}

	return nil
}

func GetResourceId(ctx context.Context, name string) (int64, error) {
	conn, err := sql.Open("sqlite", "resources.db")
	if err != nil {
		return 9999999, err
	}

	defer conn.Close()

	var id int64
	err = conn.QueryRowContext(ctx, "SELECT id FROM resourceIds WHERE name = ?", name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("resource %s not found in resourceIds", name)
	}

	return id, nil
}
