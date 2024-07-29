package dispatch

import (
	"net/http"

	queryparser "github.com/harish876/ShardPilot/lib/queryParser"
	"github.com/labstack/echo/v4"
)

type DispatchRequest struct {
	Query string `json:"query"`
}

type DispatchResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func GetDispatchHandler(c echo.Context) error {
	var body DispatchRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, DispatchResponse{Error: "Unable to unmarshal request"})
	}
	lp := queryparser.NewLogicalPlanParams([]byte(body.Query))
	lp.
		GetQueryType().
		GetShardID()

	return c.JSON(http.StatusOK, DispatchResponse{Message: lp.String()})
}

func PostDispatchHandler(c echo.Context) error {
	var body DispatchRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, DispatchResponse{Error: "Unable to unmarshal request"})
	}
	return c.JSON(http.StatusOK, "Post Dispatch")
}

func PutDispatchHandler(c echo.Context) error {
	var body DispatchRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, DispatchResponse{Error: "Unable to unmarshal request"})
	}
	return c.JSON(http.StatusOK, "Put Dispatch")
}

func DeleteDispatchHandler(c echo.Context) error {
	var body DispatchRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, DispatchResponse{Error: "Unable to unmarshal request"})
	}
	return c.JSON(http.StatusOK, "Delete Dispatch")
}
