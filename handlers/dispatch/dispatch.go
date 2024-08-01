package dispatch

import (
	"fmt"
	"net/http"

	"github.com/harish876/ShardPilot/db"
	logicalplanner "github.com/harish876/ShardPilot/lib/ast/logicalPlanner"
	physicalplanner "github.com/harish876/ShardPilot/lib/ast/physicalPlanner"
	"github.com/labstack/echo/v4"
	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type DispatchRequest struct {
	Query string `json:"query"`
}

type DispatchResponse struct {
	Message string `json:"message"`
	Data    []User `json:"data,omitempty"`
	Error   string `json:"error"`
}

type User struct {
	UserID      int
	Name        string
	PhoneNumber string
	Email       string
}

func GetDispatchHandler(c echo.Context) error {
	var body DispatchRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, DispatchResponse{Error: "Unable to unmarshal request"})
	}
	node, err := pg_query.Parse(body.Query)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			DispatchResponse{Error: fmt.Sprintf("Unabe to parse query %s", err.Error())},
		)
	}
	lp, err := logicalplanner.NewLogicalPlanParams(node)
	if err != nil {
		return c.JSON(http.StatusBadRequest, DispatchResponse{Error: err.Error()})
	}
	lp.
		GetQueryType().
		GetShardId()

	modifiedQuery, err := physicalplanner.RewriteSelectQuery(node, "shardid")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, DispatchResponse{Error: err.Error()})
	}

	conn, _ := db.GetConnectionPoolForShard(lp.ShardId)
	rows, err := conn.Query(modifiedQuery)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			DispatchResponse{Error: fmt.Sprintf("unable to query database %s", err.Error())},
		)
	}
	_ = rows
	defer rows.Close()
	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID, &user.Name, &user.PhoneNumber, &user.Email); err != nil {
			fmt.Println(err)
		}
		users = append(users, user)
	}
	return c.JSON(http.StatusOK, DispatchResponse{
		Message: fmt.Sprintf("Query on ShardID %d: %s", lp.ShardId, modifiedQuery),
		Data:    users,
	},
	)
}

func PostDispatchHandler(c echo.Context) error {
	var body DispatchRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, DispatchResponse{Error: "Unable to unmarshal request"})
	}
	return c.JSON(http.StatusOK, "Post Dispatch yet to implement")
}

func PutDispatchHandler(c echo.Context) error {
	var body DispatchRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, DispatchResponse{Error: "Unable to unmarshal request"})
	}
	return c.JSON(http.StatusOK, "Put Dispatch yet to implement")
}

func DeleteDispatchHandler(c echo.Context) error {
	var body DispatchRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, DispatchResponse{Error: "Unable to unmarshal request"})
	}
	return c.JSON(http.StatusOK, "Delete Dispatch yet to implement")
}
