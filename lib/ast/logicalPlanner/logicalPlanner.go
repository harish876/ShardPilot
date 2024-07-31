package logicalplanner

import (
	"fmt"

	"github.com/harish876/ShardPilot/lib/ast"
	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type LogicalPlanParams struct {
	err       error
	node      *pg_query.Node
	queryType string
	shardId   uint32
}

func NewLogicalPlanParams(node *pg_query.Node) *LogicalPlanParams {
	return &LogicalPlanParams{
		node: node,
	}
}

func (lp *LogicalPlanParams) HasError() bool {
	return lp.err != nil
}

func (lp *LogicalPlanParams) String() string {
	return fmt.Sprintf("Logical Shard Id: %d, Query Type: %s", lp.shardId, lp.queryType)
}

func (lp *LogicalPlanParams) GetShardID() *LogicalPlanParams {
	acc := make(map[string]int32)
	ast.GetAllColumns(lp.node.GetSelectStmt().WhereClause, acc)
	if val, ok := acc["shardid"]; ok {
		lp.shardId = uint32(val)
	}
	return lp
}

func (lp *LogicalPlanParams) GetQueryType() *LogicalPlanParams {
	qt, err := ast.GetQueryType(lp.node)
	if err != nil {
		lp.err = err
		return lp
	}
	lp.queryType = qt
	return lp
}
