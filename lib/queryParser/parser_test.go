package queryparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintAst(t *testing.T) {
	input := []byte(`Insert into users (userId,name) VALUES("a","b")`)
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
	fmt.Println(lp.err.Error())
}

func TestGetQueryType(t *testing.T) {
	input := []byte(`INSERT INTO users (userId,name) (123,"harish")`)
	lp := NewLogicalPlanParams(input)
	lp.GetQueryType()
	assert.NoError(t, lp.err, "Error getting query type from query%v", lp.err)
}
