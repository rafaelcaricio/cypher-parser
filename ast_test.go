package cypher_test

import (
	"testing"

	"github.com/rafaelcaricio/cypher-parser"
)

func TestQueryToString(t *testing.T) {
	strQuery := `MATCH (user :User {name: "Adam"}) RETURN user`

	q := cypher.Query{}
	user := cypher.Variable("user")
	node := cypher.NodePattern{
		Variable: &user,
		Labels:   []string{"User"},
		Properties: map[string]cypher.Expr{
			"name": cypher.StrLiteral("Adam"),
		},
	}
	q.Root = cypher.SingleQuery{
		Reading: []cypher.ReadingClause{
			{Pattern: []cypher.MatchPattern{
				{Elements: []cypher.PatternElement{node}},
			}},
		},
		ReturnItems: []cypher.Expr{
			user,
		},
	}

	r := q.String()
	if r != strQuery {
		t.Errorf("Did not generate correct query: \nExpected:\n\t%s\nGot:\n\t%s", strQuery, r)
	}
}
