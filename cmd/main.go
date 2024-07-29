package main

import (
	"fmt"
	"log"

	"github.com/harish876/ShardPilot/config"
	"github.com/harish876/ShardPilot/routes"
	"github.com/labstack/echo/v4"
)

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
