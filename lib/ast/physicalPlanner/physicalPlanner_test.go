package physicalplanner

import (
	"log/slog"
	"os"
	"testing"

	"github.com/harish876/ShardPilot/lib"
	"github.com/harish876/ShardPilot/lib/ast"
	pg_query "github.com/pganalyze/pg_query_go/v5"
	"github.com/stretchr/testify/assert"
)

func setup() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)
}

func TestQuerier(t *testing.T) {
	setup()
	query := "SELECT * from users where shardKey= 'user_id'"
	result, err := pg_query.Parse(query)
	if err != nil {
		panic(err)
	}
	root := result.Stmts[0]
	acc := make(map[string]interface{})
	ast.GetAllColumns(root.Stmt.GetSelectStmt().WhereClause, acc)
	_, ok := acc[lib.SHARD_KEY_IDENTIFIER]
	assert.Equal(t, true, ok)
}

func TestPgQuery1(t *testing.T) {
	setup()
	query := "SELECT * from users where shardKey = 'user_id'"
	node, err := pg_query.Parse(query)
	assert.NoError(t, err)
	modifiedQuery, err := RemoveNodeFromSelectQuery(node, lib.SHARD_KEY_IDENTIFIER)
	if err != nil {
		t.Fatalf("Error - %v", err)
	}
	assert.Equal(t, "SELECT * FROM users", modifiedQuery)
}

func TestPgQuery2(t *testing.T) {
	setup()
	query := `SELECT * from users where shardKey = 'user_id' and (userId=3 or age < 25)`
	node, err := pg_query.Parse(query)
	assert.NoError(t, err)
	modifiedQuery, err := RemoveNodeFromSelectQuery(node, lib.SHARD_KEY_IDENTIFIER)
	if err != nil {
		t.Fatalf("Error - %v", err)
	}
	assert.Equal(t, "SELECT * FROM users WHERE userid = 3 OR age < 25", modifiedQuery)
}

func TestPgQuery3(t *testing.T) {
	setup()
	query := `SELECT * from users where (userId=3 or age < 25) and shardKey = 'user_id'`
	node, err := pg_query.Parse(query)
	assert.NoError(t, err)
	modifiedQuery, err := RemoveNodeFromSelectQuery(node, lib.SHARD_KEY_IDENTIFIER)

	if err != nil {
		t.Fatalf("Error - %v", err)
	}
	assert.Equal(t, "SELECT * FROM users WHERE userid = 3 OR age < 25", modifiedQuery)
}

func TestPgQuery4(t *testing.T) {
	setup()
	query := "SELECT * from users where userId=3 or age < 25 and shardKey = 'user_id'"
	node, err := pg_query.Parse(query)
	assert.NoError(t, err)
	modifiedQuery, err := RemoveNodeFromSelectQuery(node, lib.SHARD_KEY_IDENTIFIER)
	if err != nil {
		t.Fatalf("Error - %v", err)
	}
	assert.Equal(t, "SELECT * FROM users WHERE userid = 3 OR age < 25", modifiedQuery)
}

func TestPgQuery5(t *testing.T) {
	setup()
	query := "SELECT * FROM users WHERE shardKey = 'user_id' or user_id = 123 AND age BETWEEN 25 AND 35"
	node, err := pg_query.Parse(query)
	assert.NoError(t, err)
	modifiedQuery, err := RemoveNodeFromSelectQuery(node, lib.SHARD_KEY_IDENTIFIER)
	if err != nil {
		t.Fatalf("Error - %v", err)
	}
	assert.Equal(
		t,
		"SELECT * FROM users WHERE user_id = 123 AND age BETWEEN 25 AND 35",
		modifiedQuery,
	)
}

// TODO: Debug this case
func TestPgQueryA5(t *testing.T) {
	setup()
	query := "SELECT * FROM users WHERE shardKey = 'user_id' and user_id = 123 AND age BETWEEN 25 AND 35"
	node, err := pg_query.Parse(query)
	assert.NoError(t, err)
	modifiedQuery, err := RemoveNodeFromSelectQuery(node, lib.SHARD_KEY_IDENTIFIER)
	if err != nil {
		t.Fatalf("Error - %v", err)
	}
	assert.Equal(
		t,
		"SELECT * FROM users WHERE user_id = 123 AND age BETWEEN 25 AND 35",
		modifiedQuery,
	)
}

func TestPgQuery6(t *testing.T) {
	setup()
	query := "SELECT * FROM users WHERE shardKey = 'user_id' AND name LIKE 'John%'"
	node, err := pg_query.Parse(query)
	assert.NoError(t, err)
	modifiedQuery, err := RemoveNodeFromSelectQuery(node, lib.SHARD_KEY_IDENTIFIER)
	if err != nil {
		t.Fatalf("Error - %v", err)
	}
	assert.Equal(t, "SELECT * FROM users WHERE name LIKE 'John%'", modifiedQuery)
}

func TestPgQuery7(t *testing.T) {
	setup()
	query := "SELECT * FROM users WHERE shardKey = 'user_id' AND age IN (25, 30, 35)"
	node, err := pg_query.Parse(query)
	assert.NoError(t, err)
	modifiedQuery, err := RemoveNodeFromSelectQuery(node, lib.SHARD_KEY_IDENTIFIER)
	if err != nil {
		t.Fatalf("Error - %v", err)
	}
	assert.Equal(t, "SELECT * FROM users WHERE age IN (25, 30, 35)", modifiedQuery)
}

func TestPgQuery8(t *testing.T) {
	setup()
	query := "SELECT * FROM users WHERE shardKey = 'user_id' AND email IS NOT NULL;"
	node, err := pg_query.Parse(query)
	assert.NoError(t, err)
	modifiedQuery, err := RemoveNodeFromSelectQuery(node, lib.SHARD_KEY_IDENTIFIER)
	if err != nil {
		t.Fatalf("Error - %v", err)
	}
	assert.Equal(t, "SELECT * FROM users WHERE email IS NOT NULL", modifiedQuery)
}

func TestPgQuery9(t *testing.T) {
	setup()
	query := "SELECT * FROM users INNER JOIN favourites on users.id = favourites.user_id where shardKey = 'user_id'"
	node, err := pg_query.Parse(query)
	assert.NoError(t, err)
	modifiedQuery, err := RemoveNodeFromSelectQuery(node, lib.SHARD_KEY_IDENTIFIER)
	if err != nil {
		t.Fatalf("Error - %v", err)
	}
	assert.Equal(
		t,
		"SELECT * FROM users JOIN favourites ON users.id = favourites.user_id",
		modifiedQuery,
	)
}
