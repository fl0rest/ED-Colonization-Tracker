package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/logging"
	"api-gateway/internal/middleware"
	"context"
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

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/api/fetch", handlers.FetchHandler)
	mux.HandleFunc("POST /api/save", handlers.SaveHandler)

	loggedMux := middleware.LoggingMiddleware(mux)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: loggedMux,
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Sys("Graceful shutdown failed:", err)
		log.Sys("Forcing Shutdown")
		os.Exit(1)
	}

	log.Sys("Server shutdown")
}
