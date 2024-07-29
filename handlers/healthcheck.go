package handlers

import (
	"fmt"
	"net/http"

	"github.com/harish876/ShardPilot/config"
	"github.com/labstack/echo/v4"
)

func HealthcheckHandler(c echo.Context) error {
	cfg, _ := config.GetConfig()
	return c.JSON(http.StatusOK, fmt.Sprintf("Server running on port %d", cfg.Server))
}
