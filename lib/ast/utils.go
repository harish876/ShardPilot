package ast

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

type QueryExecutionParams struct {
	Cursor     *sitter.QueryCursor
	Query      *sitter.Query
	Node       *sitter.Node
	sourceCode []byte
}

func NewQueryExecutionParams(
	cursor *sitter.QueryCursor,
	query *sitter.Query,
	node *sitter.Node,
	sourceCode []byte,
) *QueryExecutionParams {
	return &QueryExecutionParams{
		Cursor:     cursor,
		Query:      query,
		Node:       node,
		sourceCode: sourceCode,
	}
}

func NewQueryCursor(
	lang *sitter.Language,
	sourceCode []byte,
	query []byte,
) (*QueryExecutionParams, error) {
	node, _ := sitter.ParseCtx(context.Background(), sourceCode, lang)

	sitterQuery, err := sitter.NewQuery(query, lang)
	if err != nil {
		return nil, err
	}
	queryCursor := sitter.NewQueryCursor()

	return NewQueryExecutionParams(queryCursor, sitterQuery, node, sourceCode), nil
}

func (qc *QueryExecutionParams) Exec() {
	qc.Cursor.Exec(qc.Query, qc.Node)
}

func (qc *QueryExecutionParams) NextMatch() (*sitter.QueryMatch, bool) {
	return qc.Cursor.NextMatch()
}

func (qc *QueryExecutionParams) GetCaptures(m *sitter.QueryMatch) *sitter.QueryMatch {
	return qc.Cursor.FilterPredicates(m, qc.sourceCode)
}
