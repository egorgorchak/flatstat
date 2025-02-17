package main

import (
	"context"
	"flatstat/internal/config"
	"flatstat/internal/handlers"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	l := setupLogger(cfg.Env)

	l.Info("Starting flatstat", slog.String("env", cfg.Env))
	l.Debug("Debug is enabled")

	sm := http.NewServeMux()

	h := handlers.Info{}

	sm.Handle("/info", &h)

	s := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      sm,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}

	go func() {
		l.Info("Starting server", "address", cfg.HTTPServer.Address)

		err := s.ListenAndServe()
		if err != nil {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Info("Received terminate, graceful shutdown", "signal", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
