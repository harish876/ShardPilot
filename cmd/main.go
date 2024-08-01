package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/harish876/ShardPilot/config"
	"github.com/harish876/ShardPilot/db"
	"github.com/harish876/ShardPilot/routes"
	"github.com/labstack/echo/v4"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	slog.Info("init", "init slog", "Structured Logging Setup")

	if err := db.Init(); err != nil {
		slog.Error("db.Init", "failed to init database connections", err)
		os.Exit(1)
	}
}

func main() {
	app := echo.New()
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Unable to InitConfig %v", err)
	}

	routes.RegisterHealthCheckRoutes(app)
	routes.RegisterDispatchRoutes(app)
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", cfg.Server)))
}
