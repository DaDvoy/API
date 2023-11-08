package main

import (
	postgres "API/internal/DB"
	server "API/internal/app"
	"API/internal/config"
	"API/internal/lib/logger/sl"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
	"log"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func main() {
	cfg := config.MustLoad()

	logs := setupLogger(cfg.Env)
	logs.Info("starting", slog.String("env", cfg.Env))
	logs.Debug("debug messages are enabled")
	db, err := postgres.New()
	if err != nil {
		logs.Error("Failed to init database", sl.Err(err)) // TODO: If I want not use sl, need to delete it
		os.Exit(1)
	}
	if err := postgres.CreateTable(db); err != nil {
		log.Fatal(err)
	}
	s := server.New(logs, cfg)
	if err := s.Start(db); err != nil {
		log.Fatal(err)
	}
	defer postgres.CloseDB(db)
}
