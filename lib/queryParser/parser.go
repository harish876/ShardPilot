package queryparser

import (
	"github.com/smacker/go-tree-sitter/sql"
)

func GetShardID(input []byte) ([]string, error) {
	result := make([]string, 0)
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
	qc, err := NewQueryCursor(lang, input, query)

	if err != nil {
		return nil, err
	}

	if qc.node.HasError() {
		return nil, err
	}

	qc.exec()
	for {
		m, ok := qc.nextMatch()
		if !ok {
			break
		}
		m = qc.getCaptures(m)
		//fmt.Println(m.Captures)
		for _, c := range m.Captures {
			//fmt.Println("Capture group", c.Node.Type(), c.Node.Content(input), c.Node.IsMissing())
			if c.Node.Type() == "literal" {
				result = append(result, c.Node.Content(input))
			}
		}
	}
	return result, nil
}
