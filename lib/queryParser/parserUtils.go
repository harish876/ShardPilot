package queryparser

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

type QueryExecutionParams struct {
	cursor     *sitter.QueryCursor
	query      *sitter.Query
	node       *sitter.Node
	sourceCode []byte
}

func NewQueryExecutionParams(
	cursor *sitter.QueryCursor,
	query *sitter.Query,
	node *sitter.Node,
	sourceCode []byte,
) *QueryExecutionParams {
	return &QueryExecutionParams{
		cursor:     cursor,
		query:      query,
		node:       node,
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

func (qc *QueryExecutionParams) exec() {
	qc.cursor.Exec(qc.query, qc.node)
}

func (qc *QueryExecutionParams) nextMatch() (*sitter.QueryMatch, bool) {
	return qc.cursor.NextMatch()
}

func (qc *QueryExecutionParams) getCaptures(m *sitter.QueryMatch) *sitter.QueryMatch {
	return qc.cursor.FilterPredicates(m, qc.sourceCode)
}
