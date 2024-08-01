package routes

import (
	"github.com/harish876/ShardPilot/app/handlers/dispatch"
	"github.com/labstack/echo/v4"
)

func RegisterDispatchRoutes(e *echo.Echo) {
	rg := e.Group("")
	rg.GET(
		"/dispatch",
		dispatch.GetDispatchHandler,
	)

	rg.POST(
		"/dispatch",
		dispatch.PostDispatchHandler,
	)

	rg.PUT(
		"/dispatch",
		dispatch.PutDispatchHandler,
	)

	rg.DELETE(
		"/dispatch",
		dispatch.DeleteDispatchHandler,
	)
}
