package queryparser

import (
	"context"
	"fmt"
	"strconv"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/sql"
)

type LogicalPlanParams struct {
	err              error
	queryType        string
	input            []byte
	shardId          uint32
	isShardIdPresent bool
}

func NewLogicalPlanParams(input []byte) *LogicalPlanParams {
	return &LogicalPlanParams{
		input: input,
		err:   nil,
	}
}

func (lp *LogicalPlanParams) PrintAst() (string, error) {
	n, err := sitter.ParseCtx(context.Background(), lp.input, sql.GetLanguage())
	if err != nil {
		return "", err
	}
	return n.String(), nil
}

func (lp *LogicalPlanParams) HasError() bool {
	return lp.err != nil
}

func (lp *LogicalPlanParams) String() string {
	if lp.isShardIdPresent {
		return fmt.Sprintf("Logical Shard Id: %d, Query Type: %s", lp.shardId, lp.queryType)
	} else {
		return fmt.Sprintf("Logical Shard Id: %s, Query Type: %s", "NA", lp.queryType)

	}
}

func (lp *LogicalPlanParams) GetShardID() *LogicalPlanParams {
	var result string
	lang := sql.GetLanguage()
	//TODO: make query better
	query := []byte(
		`(
       (where
          predicate: (
            [
              (binary_expression
                (binary_expression)*
                left: (field name: (identifier) @id)
                right: (literal) @shardId
                (#eq? @id "shardId")
              )


              (binary_expression(binary_expression
                left: (field name: (identifier) @id)
                right: (literal) @shardId
                (#eq? @id "shardId")
              ))            
            ]
        )
      ) 
    )`,
	)
	qc, err := NewQueryCursor(lang, lp.input, query)
	lp.err = err

	if qc.node.HasError() {
		lp.err = err
	}

	qc.exec()
	for {
		m, ok := qc.nextMatch()
		if !ok {
			break
		}
		m = qc.getCaptures(m)
		for _, c := range m.Captures {
			if c.Node.Type() == "literal" {
				result = c.Node.Content(lp.input)
				break
			}
		}
	}

	shardId, err := strconv.Atoi(result)
	if err != nil {
		lp.err = fmt.Errorf("incorrect shardId, not a valid uint32: %s", result)
	}
	lp.shardId = uint32(shardId)
	lp.isShardIdPresent = (err == nil)
	return lp
}

func (lp *LogicalPlanParams) GetQueryType() *LogicalPlanParams {
	lang := sql.GetLanguage()
	query := []byte(
		`(
      (program
          (statement
            [ 
              (
                insert (keyword_insert) @queryType
              )

              (
                 select (keyword_select) @queryType
              )

              (
                 update (keyword_update) @queryType
              )

              (
                 delete (keyword_delete) @queryType
              )            
            ]
          )
      ) 
    )
`,
	)
	qc, err := NewQueryCursor(lang, lp.input, query)
	lp.err = err

	if qc.node.HasError() {
		lp.err = err
	}

	var queryType string
	qc.exec()
	for {
		m, ok := qc.nextMatch()
		if !ok {
			break
		}
		m = qc.getCaptures(m)
		for _, c := range m.Captures {
			switch c.Node.Type() {
			case "keyword_select":
				queryType = "SELECT"
			case "keyword_update":
				queryType = "UPDATE"
			case "keyword_insert":
				queryType = "INSERT"
			case "keyword_delete":
				queryType = "DELETE"
			default:
				queryType = "UNKNOWN"
				lp.err = fmt.Errorf("unknown query type")
			}
		}
	}
	lp.queryType = queryType
	return lp
}

func (lp *LogicalPlanParams) Validate() error {
	return nil
}
