package logicalplanner

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetShardID1(t *testing.T) {
	input := `select * from users where shardKey = 'user_id' and user_id = 123`
	root, err := pg_query.Parse(input)
	if err != nil {
		t.Fatal("error building AST", err)
	}
	lp, err := NewLogicalPlanParams(root)
	assert.NoError(t, err)
	lp.
		GetQueryType().
		GetShardId()

	assert.NoError(t, lp.err, "Error creating query plan")
	assert.Equal(t, lp.ShardId, uint32(123))
	assert.Equal(t, lp.QueryType, "SELECT")
}

func TestGetShardID2(t *testing.T) {
	input := `select * from users where shardKey = 'user_id' and user_id = 123`
	root, err := pg_query.Parse(input)
	if err != nil {
		t.Fatal("error building AST", err)
	}
	lp, err := NewLogicalPlanParams(root)
	assert.NoError(t, err)
	lp.
		GetQueryType().
		GetShardId()

	assert.NoError(t, lp.err, "Error creating query plan")
	assert.Equal(t, lp.ShardId, uint32(123))
	assert.Equal(t, lp.QueryType, "SELECT")
}

func TestGetShardID3(t *testing.T) {
	input := `select * from users`
	root, err := pg_query.Parse(input)
	if err != nil {
		t.Fatal("error building AST", err)
	}

	lp, err := NewLogicalPlanParams(root)
	assert.NoError(t, err)
	lp.
		GetQueryType().
		GetShardId()

	assert.Equal(t, lp.ShardId, uint32(0))
	assert.Equal(t, lp.QueryType, "SELECT")
}

func TestGetShardID4(t *testing.T) {
	input := `select * from users where shardKey = 'user_id' and colX = Z`
	root, err := pg_query.Parse(input)
	if err != nil {
		t.Fatal("error building AST", err)
	}
	lp, err := NewLogicalPlanParams(root)
	assert.NoError(t, err)
	lp.
		GetQueryType().
		GetShardId()
	assert.Equal(t, lp.ShardId, uint32(0))
	assert.Equal(t, lp.QueryType, "SELECT")
}

func TestGetQueryType(t *testing.T) {
	input := `INSERT INTO employees (id, name, position) VALUES (1, 'John Doe', 'Software Engineer');`
	root, err := pg_query.Parse(input)
	if err != nil {
		t.Fatal("error building AST", err)
	}
	lp, err := NewLogicalPlanParams(root)
	assert.NoError(t, err)
	lp.GetQueryType()
	assert.NoError(t, lp.err, "Error getting query type from query%v", lp.err)
}
