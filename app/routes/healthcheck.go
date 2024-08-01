package routes

import (
	"github.com/harish876/ShardPilot/app/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterHealthCheckRoutes(e *echo.Echo) {
	rg := e.Group("")
	rg.GET("/healthcheck", handlers.HealthcheckHandler)
}
