package logicalplanner

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/harish876/ShardPilot/lib"
	"github.com/harish876/ShardPilot/lib/ast"
	"github.com/harish876/ShardPilot/lib/hash"
	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type LogicalPlanParams struct {
	err       error
	node      *pg_query.Node
	colMap    map[string]interface{}
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
		node:   node,
		colMap: make(map[string]interface{}),
	}, nil
}

func (lp *LogicalPlanParams) HasError() bool {
	return lp.err != nil
}

func (lp *LogicalPlanParams) GetShardId() *LogicalPlanParams {
	if lp.node == nil || lp.QueryType != "SELECT" {
		return lp
	}
	shardKey, ok := lp.collectColNamesFromWhereClause().checkIfShardKeyIsPresent()
	if !ok {
		slog.Info("GetShardId", "Shard key is not present in the query", ok)
		return lp
	}

	if val, ok := lp.getShardKeyValue(shardKey); !ok {
		return lp
	} else {
		shardId, _ := hash.CalculateShardId(hash.IntToBytes(int(val)), 3) //todo hard coded
		lp.ShardId = shardId
	}
	return lp
}

func (lp *LogicalPlanParams) getShardKeyValue(shardKey string) (int32, bool) {
	if _, ok := lp.colMap[shardKey]; !ok {
		return -1, false
	}
	val := lp.colMap[shardKey]
	switch reflect.ValueOf(val).Kind() {
	case reflect.Int32:
		return val.(int32), true
	default:
		return -1, false
	}
}

func (lp *LogicalPlanParams) collectColNamesFromWhereClause() *LogicalPlanParams {
	if lp.HasError() {
		return lp
	}
	ast.GetAllColumns(lp.node.GetSelectStmt().WhereClause, lp.colMap)
	return lp
}

func (lp *LogicalPlanParams) checkIfShardKeyIsPresent() (string, bool) {
	if val, ok := lp.colMap[lib.SHARD_KEY_IDENTIFIER]; ok {
		if reflect.ValueOf(val).Kind() == reflect.String {
			return val.(string), true
		} else {
			lp.err = fmt.Errorf("invalid shard key. shard key has to be a string %v", val)
			return "", false
		}
	} else {
		return "", false
	}
}

func (lp *LogicalPlanParams) GetQueryType() *LogicalPlanParams {
	if lp.HasError() {
		return lp
	}

	qt, err := ast.GetQueryType(lp.node)
	if err != nil {
		lp.err = err
		return lp
	}
	lp.QueryType = qt
	return lp
}
