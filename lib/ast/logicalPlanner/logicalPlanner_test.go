package logicalplanner

import (
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetShardID1(t *testing.T) {
	input := `select * from users where shardId =1`
	root, err := pg_query.Parse(input)
	if err != nil {
		t.Fatal("error building AST", err)
	}
	node := root.Stmts[0].Stmt
	lp := NewLogicalPlanParams(node)
	lp.
		GetQueryType().
		GetShardID()

	assert.NoError(t, lp.err, "Error creating query plan")
	assert.Equal(t, lp.shardId, uint32(1))
	assert.Equal(t, lp.queryType, "SELECT")
}

func TestGetShardID2(t *testing.T) {
	input := `select * from users where shardId = 83310 and userId = 123`
	root, err := pg_query.Parse(input)
	if err != nil {
		t.Fatal("error building AST", err)
	}
	node := root.Stmts[0].Stmt
	lp := NewLogicalPlanParams(node)
	lp.
		GetQueryType().
		GetShardID()

	assert.NoError(t, lp.err, "Error creating query plan")
	assert.Equal(t, lp.shardId, uint32(83310))
	assert.Equal(t, lp.queryType, "SELECT")
}

func TestGetShardID3(t *testing.T) {
	input := `select * from users`
	root, err := pg_query.Parse(input)
	if err != nil {
		t.Fatal("error building AST", err)
	}

	node := root.Stmts[0].Stmt
	lp := NewLogicalPlanParams(node)
	lp.
		GetQueryType().
		GetShardID()

	assert.Equal(t, lp.shardId, uint32(0))
	assert.Equal(t, lp.queryType, "SELECT")
}

func TestGetShardID4(t *testing.T) {
	input := `select * from users where shardId=1 and (userId=3 or colX = Z)`
	root, err := pg_query.Parse(input)
	if err != nil {
		t.Fatal("error building AST", err)
	}
	node := root.Stmts[0].Stmt
	lp := NewLogicalPlanParams(node)
	lp.
		GetQueryType().
		GetShardID()

	assert.Equal(t, lp.shardId, uint32(1))
	assert.Equal(t, lp.queryType, "SELECT")
}

func TestGetQueryType(t *testing.T) {
	input := `INSERT INTO employees (id, name, position) VALUES (1, 'John Doe', 'Software Engineer');`
	root, err := pg_query.Parse(input)
	if err != nil {
		t.Fatal("error building AST", err)
	}
	node := root.Stmts[0].Stmt
	lp := NewLogicalPlanParams(node)
	lp.GetQueryType()
	assert.NoError(t, lp.err, "Error getting query type from query%v", lp.err)
}
