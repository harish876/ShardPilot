package ast

import (
	"fmt"
	"log/slog"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func IsAExpr(n *pg_query.Node) bool {
	return n.GetAExpr() != nil
}

func IsBoolExpr(n *pg_query.Node) bool {
	return n.GetBoolExpr() != nil
}

func GetColumnNameFromLexprNode(node *pg_query.A_Expr) (string, error) {
	if node == nil {
		return "", fmt.Errorf("nil node")
	} else if node.Lexpr.GetColumnRef() == nil {
		return "", fmt.Errorf("no columns referenced %s", node.Lexpr.String())
	} else if len(node.Lexpr.GetColumnRef().Fields) == 0 {
		return "", fmt.Errorf("no columns referenced %s", node.Lexpr.String())
	} else {
		return node.Lexpr.GetColumnRef().Fields[0].GetString_().Sval, nil
	}
}

func GetAllColumns(node *pg_query.Node, acc map[string]interface{}) {
	if node == nil {
		return
	}
	if IsAExpr(node) && node.GetAExpr().Name[0].GetString_().Sval == "=" {
		colName, _ := GetColumnNameFromLexprNode(node.GetAExpr())
		value := node.GetAExpr().GetRexpr().GetAConst()
		if value.GetIval() != nil {
			acc[colName] = value.GetIval().Ival
		} else if value.GetSval() != nil {
			acc[colName] = value.GetSval().Sval
		}
	} else if IsBoolExpr(node) {
		parent := node.GetBoolExpr()
		left := parent.Args[0]
		right := parent.Args[1]
		GetAllColumns(left, acc)
		GetAllColumns(right, acc)
	} else {
		slog.Debug("GetAllColumns", "not recognised", node.String())
	}
}

func GetQueryType(node *pg_query.Node) (string, error) {
	if node != nil {
		return "SELECT", nil
	} else if node.GetInsertStmt() != nil {
		return "INSERT", nil
	} else if node.GetUpdateStmt() != nil {
		return "UPDATE", nil
	} else if node.GetDeleteStmt() != nil {
		return "DELETE", nil
	} else {
		return "", fmt.Errorf("unknown query type %v", node)
	}
}

func GetTableNameFromSelect(node *pg_query.Node) (string, error) {
	return node.String(), nil
}
