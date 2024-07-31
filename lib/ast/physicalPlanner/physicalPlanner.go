package physicalplanner

import (
	"fmt"

	"log/slog"

	"github.com/harish876/ShardPilot/lib/ast"
	pg_query "github.com/pganalyze/pg_query_go/v5"
)

// TODO: improve code structure
func RewriteSelectQuery(query string) (string, error) {
	result, err := pg_query.Parse(query)
	if err != nil {
		panic(err)
	}

	if len(result.Stmts) == 0 {
		return "", fmt.Errorf("no statements in the query - %s", query)
	}

	root := result.Stmts[0]
	whereClause := root.Stmt.GetSelectStmt().WhereClause

	if ast.IsBoolExpr(whereClause) {
		slog.Debug("Rewrite Query", "Query", query)
		slog.Debug("RewriteQuery", "Original Where Clause", whereClause.String())
		modifiedWhereClause := removeShardId(whereClause, "shardid")
		slog.Debug("Rewrite Query", "Modified Where Clause", modifiedWhereClause.String())
		root.Stmt.GetSelectStmt().WhereClause = modifiedWhereClause

	} else if ast.IsAExpr(whereClause) {
		colName, err := ast.GetColumnNameFromLexprNode(whereClause.GetAExpr())
		if err != nil {
			return "", fmt.Errorf("unable to get column names from where clause %v", err)
		}
		if colName == "shardid" {
			root.Stmt.GetSelectStmt().WhereClause = nil
		}
	} else {
		slog.Debug("RewriteQuery", "unhandled case", whereClause)
		return query, fmt.Errorf("unhandled case")
	}

	modifiedQuery, err := pg_query.Deparse(result)
	if err != nil {
		return "", fmt.Errorf("error deparsing query - %v", err)
	}
	return modifiedQuery, nil
}

func removeShardId(node *pg_query.Node, key string) *pg_query.Node {
	if node == nil {
		return nil
	} else if ast.IsBoolExpr(node) {
		binaryExpr := node.GetBoolExpr()
		left := binaryExpr.Args[0]
		right := binaryExpr.Args[1]

		colNameMatchesLeft, _ := ast.GetColumnNameFromLexprNode(left.GetAExpr())
		colNameMatchesRight, _ := ast.GetColumnNameFromLexprNode(right.GetAExpr())

		slog.Debug("removeShardId", "Old Right Child -- ", right.String()+"\n")
		slog.Debug("removeShardId", "Old Left Child --", left.String()+"\n")
		slog.Debug("removeShardId", "Old Parent --", node.String()+"\n")

		if colNameMatchesLeft == key {
			node = right
			left = nil
		} else if colNameMatchesRight == key {
			node = left
			right = nil
		}

		slog.Debug("removeShardId", "New Parent --", node.String()+"\n")
		binaryExpr.Args[0] = removeShardId(left, key)
		binaryExpr.Args[1] = removeShardId(right, key)
	}
	return node
}
