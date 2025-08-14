package main

import (
	"context"
	"ed-tracker/internal/config"
	"ed-tracker/internal/db"
	"ed-tracker/internal/handlers"
	"ed-tracker/internal/logging"
	"ed-tracker/internal/middleware"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := logging.Log
	cfg := config.Load()
	mux := http.NewServeMux()
	ctx := context.Background()

	factory := &db.DBFactory{Path: cfg.DbFile}
	queries, sqlDb, err := factory.Connect(ctx)
	if err != nil {
		log.Error("Failed to connect to db:", err)
	}
	defer sqlDb.Close()

	handlers.Init(queries)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Sys("Running ParseLatestEvents")
				if err := handlers.ParseLatestEvent(ctx, queries); err != nil {
					log.Error("ParseLatestEvents Error:", err)
				}
			}
		}
	}()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("POST /api/save", handlers.SaveHandler)
	mux.HandleFunc("/events", handlers.SseHandler)
	mux.HandleFunc("/static/", handlers.StaticHandler)

	loggedMux := middleware.LoggingMiddleware(mux)

	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	defer shutdownCancel()

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: loggedMux,
		BaseContext: func(net.Listener) context.Context {
			return shutdownCtx
		},
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		log.Sys("Starting server on", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed to start: ", err)
		}
	}()

	<-quit
	log.Sys("Shutdown Initiated")

	shutdownCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Sys("Graceful shutdown failed:", err)
		log.Sys("Forcing Shutdown")
		os.Exit(1)
	}

	log.Sys("Server shutdown")
}
