package cypher_test

import (
	"strings"
	"testing"

	"github.com/rafaelcaricio/cypher-parser"
)

func TestParseNodePatterns(t *testing.T) {
	for _, query := range []struct {
		in  string
		out string
	}{
		{
			in:  "MATCH (p) RETURN",
			out: "MATCH (p) RETURN",
		},
		{
			in:  "MATCH ((((())))) RETURN",
			out: "MATCH () RETURN",
		},
		{
			in:  "MATCH () RETURN",
			out: "MATCH () RETURN",
		},
		{
			in:  "MATCH (p :Person) RETURN",
			out: "MATCH (p :Person) RETURN",
		},
		{
			in:  "MATCH ((p :Person)) RETURN",
			out: "MATCH (p :Person) RETURN",
		},
		{
			in:  "MATCH (p :Person :Human) RETURN",
			out: "MATCH (p :Person :Human) RETURN",
		},
		{
			in:  "MATCH ( :Human) RETURN",
			out: "MATCH ( :Human) RETURN",
		},
	} {
		q, err := cypher.ParseQuery(query.in)
		if err != nil {
			t.Errorf("%s", err)
		}
		if strings.Trim(q.String(), " ") != query.out {
			t.Errorf("\nExpected:\n\t%s\nGot:\n\t%s", query.out, q)
		}
	}
}
