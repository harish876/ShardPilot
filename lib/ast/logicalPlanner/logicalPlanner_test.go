package logicalplanner

import (
	"context"
	"fmt"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/sql"
	"github.com/stretchr/testify/assert"
)

func TestPrintAst(t *testing.T) {
	input := []byte(`select * from users where shardId=1 and userId=2 or age < 25`)
	lp := NewLogicalPlanParams(input)
	ast, err := lp.PrintAst()
	assert.NoError(t, err)
	fmt.Println(ast)
}

func TestGetShardID1(t *testing.T) {
	input := []byte(`select * from users where shardId = 83310`)
	lp := NewLogicalPlanParams(input)
	lp.
		GetQueryType().
		GetShardID()

	assert.NoError(t, lp.err, "Error creating query plan")
	assert.Equal(t, lp.shardId, uint32(83310))
	assert.Equal(t, lp.queryType, "SELECT")
}

func TestGetShardID2(t *testing.T) {
	input := []byte(`select * from users where shardId = 83310 and userId = 123`)
	lp := NewLogicalPlanParams(input)
	lp.
		GetQueryType().
		GetShardID()

	assert.NoError(t, lp.err, "Error creating query plan")
	assert.Equal(t, lp.shardId, uint32(83310))
	assert.Equal(t, lp.queryType, "SELECT")
}

func TestGetShardID3(t *testing.T) {
	input := []byte(`select * from users`)
	lp := NewLogicalPlanParams(input)
	lp.
		GetQueryType().
		GetShardID()

	assert.Equal(t, lp.shardId, uint32(0))
	assert.Equal(t, lp.queryType, "SELECT")
}

func TestGetShardID4(t *testing.T) {
	input := []byte(`select * from users where shardId=1 and (userId=3 or colX = Z)`)
	lp := NewLogicalPlanParams(input)
	lp.
		GetQueryType().
		GetShardID()

	assert.Equal(t, lp.shardId, uint32(1))
	assert.Equal(t, lp.queryType, "SELECT")
}

func TestGetQueryType(t *testing.T) {
	input := []byte(`INSERT INTO users (userId,name) (123,"harish")`)
	lp := NewLogicalPlanParams(input)
	lp.GetQueryType()
	assert.NoError(t, lp.err, "Error getting query type from query%v", lp.err)
}

func TestWalk(t *testing.T) {
	input := []byte(`select * from users where shardId = 83310 and userId = 123`)
	node, _ := sitter.ParseCtx(context.Background(), input, sql.GetLanguage())
	Walk(node, input, 0)
	fmt.Println(node.String())
}
