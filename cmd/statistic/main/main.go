package main

import (
	"gorm.io/gorm"
	"log/slog"
	"os"
	"statistic-service/config"
	"statistic-service/internal/domain/statistic/delievery/http"
	"statistic-service/internal/domain/statistic/delievery/http/handlers"
	"statistic-service/internal/infrastructure/database/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	var dbClient *gorm.DB
	cfg := config.MustLoadConfig("config/example.config.yaml")
	log := setupLogger(cfg.Env)
	log.Info("Logger started successfully")
	pgConnect := postgres.NewPostgresConnect(cfg)
	client, err := pgConnect.Connect()
	if err != nil {
		panic(err)
	}
	dbClient = client.(*gorm.DB)

	defer func(pgConnect postgres.PostgresConnectable, i interface{}) {
		err := pgConnect.CloseConnection(i)
		if err != nil {
			panic(err)
		}
	}(pgConnect, dbClient)
	statisticHandler := handlers.NewStatisticHandlers(cfg, log, dbClient)
	httpServer := http.NewHTTPServer(cfg, log, statisticHandler)
	httpServer.Run()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev, envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
