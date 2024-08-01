package logicalplanner

import (
	"fmt"
	"log/slog"

	"github.com/harish876/ShardPilot/lib/ast"
	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type LogicalPlanParams struct {
	err       error
	node      *pg_query.Node
	QueryType string
	ShardId   uint32
}

func NewLogicalPlanParams(ast *pg_query.ParseResult) (*LogicalPlanParams, error) {
	stmts := ast.Stmts
	var node *pg_query.Node
	if len(stmts) == 0 {
		slog.Info("NewLogicalPlanParams", "invalid sql statement", "0 statements present")
		return nil, fmt.Errorf("invalid SQL statement. only supports single line SQL stmts")
	}
	node = stmts[0].Stmt
	return &LogicalPlanParams{
		node: node,
	}, nil
}

func (lp *LogicalPlanParams) HasError() bool {
	return lp.err != nil
}

func (lp *LogicalPlanParams) String() string {
	return fmt.Sprintf("Logical Shard Id: %d, Query Type: %s", lp.ShardId, lp.QueryType)
}

func (lp *LogicalPlanParams) GetShardId() *LogicalPlanParams {
	if lp.node == nil {
		return lp
	}

	acc := make(map[string]interface{})
	ast.GetAllColumns(lp.node.GetSelectStmt().WhereClause, acc)
	if val, ok := acc["shardid"]; ok {
		lp.ShardId = uint32(val.(int32))
	}
	return lp
}

func (lp *LogicalPlanParams) GetQueryType() *LogicalPlanParams {
	qt, err := ast.GetQueryType(lp.node)
	if err != nil {
		lp.err = err
		return lp
	}
	lp.QueryType = qt
	return lp
}
