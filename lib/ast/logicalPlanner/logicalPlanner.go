package logicalplanner

import (
	"context"
	"fmt"
	"strconv"

	"github.com/harish876/ShardPilot/lib/ast"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/sql"
)

type LogicalPlanParams struct {
	err       error
	queryType string
	input     []byte
	shardId   uint32
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
	return fmt.Sprintf("Logical Shard Id: %d, Query Type: %s", lp.shardId, lp.queryType)
}

func (lp *LogicalPlanParams) GetShardID() *LogicalPlanParams {
	var result string
	lang := sql.GetLanguage()
	//TODO: make query better
	query := []byte(
		`(
       (where 
          predicate: (_) @predicate
       )  
     )`,
	)
	qc, err := ast.NewQueryCursor(lang, lp.input, query)
	lp.err = err

	if qc.Node.HasError() {
		lp.err = err
	}

	qc.Exec()
	acc := make([]SelectLeafNode, 0)
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		m = qc.GetCaptures(m)
		for _, c := range m.Captures {
			resconstructedQuery, err := ReconstructCondition(c.Node, lp.input)
			if err != nil {
				fmt.Println("ERROR", err.Error())
			}
			fmt.Println(resconstructedQuery)
			collectLeaves(c.Node, lp.input, &acc)
		}
	}

	for _, leaf := range acc {
		if leaf.Field == "shardId" {
			result = leaf.Value
			break
		}
	}

	fmt.Println("Captures", acc)

	shardId, err := strconv.Atoi(result)
	if err != nil {
		lp.err = fmt.Errorf("incorrect shardId, not a valid uint32: %s", result)
	}
	lp.shardId = uint32(shardId)
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
	qc, err := ast.NewQueryCursor(lang, lp.input, query)
	lp.err = err

	if qc.Node.HasError() {
		lp.err = err
	}

	var queryType string
	qc.Exec()
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		m = qc.GetCaptures(m)
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

func Walk(node *sitter.Node, input []byte, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += " "
	}

	for idx := 0; idx < int(node.NamedChildCount()); idx++ {
		child := node.NamedChild(idx)

		fmt.Printf(
			"%s%s: %s %s\n",
			indent,
			node.Type(),
			node.Content(input),
			node.FieldNameForChild(idx),
		)
		Walk(child, input, depth+1)
	}
}

type SelectLeafNode struct {
	Field string
	Value string
}

func collectLeaves(node *sitter.Node, source []byte, acc *[]SelectLeafNode) {
	if node == nil {
		return
	}

	if node.Type() != "binary_expression" {
		return
	}

	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")

	var leafNode SelectLeafNode
	leafNode.Field = left.Content(source)
	leafNode.Value = right.Content(source)
	*acc = append(*acc, leafNode)

	collectLeaves(left, source, acc)
	collectLeaves(right, source, acc)
}

func ReconstructCondition(node *sitter.Node, source []byte) (string, error) {
	switch node.Type() {
	case "binary_expression":
		leftNode := node.Child(0)
		rightNode := node.Child(1)
		operator := node.Child(2).Type() // Assuming the operator is at the third position

		leftExpr, err := ReconstructCondition(leftNode, source)
		if err != nil {
			return "", err
		}
		rightExpr, err := ReconstructCondition(rightNode, source)
		if err != nil {
			return "", err
		}

		switch operator {
		case "keyword_and":
			return fmt.Sprintf("(%s AND %s)", leftExpr, rightExpr), nil
		case "keyword_or":
			return fmt.Sprintf("(%s OR %s)", leftExpr, rightExpr), nil
		case "=":
			return fmt.Sprintf("%s = %s", leftExpr, rightExpr), nil
		default:
			fmt.Println("Default Condition at operator", node.Content(source), operator)
			return node.Content(source), nil
			//return "", fmt.Errorf("unknown operator: %s", operator)
		}

	case "field":
		return node.Content(source), nil

	case "literal":
		return fmt.Sprintf("'%s'", node.Content(source)), nil

	case "=":
		return node.Content(source), nil

	default:
		fmt.Println("Default Condition at Node", node.Content(source), node.Type())
		return node.Content(source), nil
	}
}
