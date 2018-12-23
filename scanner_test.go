package cypher_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/rafaelcaricio/cypher-parser"
)

func TestScanInput(t *testing.T) {
	for _, tc := range []struct {
		in  string
		tok cypher.Token
		lit string
	}{
		{in: `something`, tok: cypher.IDENT, lit: "something"},
		{in: `desc`, tok: cypher.DESC, lit: ""},
		{in: `match`, tok: cypher.MATCH, lit: ""},
		{in: `or`, tok: cypher.OR, lit: ""},
		{in: `1233`, tok: cypher.INTEGER, lit: "1233"},
		{in: `3.14`, tok: cypher.NUMBER, lit: "3.14"},
		{in: `true`, tok: cypher.TRUE, lit: ""},
		{in: `null`, tok: cypher.NULL, lit: ""},
		{in: `"Hello, world!"`, tok: cypher.STRING, lit: "Hello, world!"},
		{in: `"String\nwith\nnewline"`, tok: cypher.STRING, lit: "String\nwith\nnewline"},
		{in: `'String\n'`, tok: cypher.STRING, lit: "String\n"},
		{in: ``, tok: cypher.EOF, lit: ""},
		{in: `<>`, tok: cypher.NEQ, lit: ""},
		{in: `<`, tok: cypher.LT, lit: ""},
		{in: `<=1`, tok: cypher.LTE, lit: ""},
		{in: `>`, tok: cypher.GT, lit: ""},
		{in: `>=`, tok: cypher.GTE, lit: ""},
		{in: `..`, tok: cypher.DOUBLEDOT, lit: ""},
		{in: `+`, tok: cypher.PLUS, lit: ""},
		{in: `+=`, tok: cypher.INC, lit: ""},
		{in: `//nice try`, tok: cypher.COMMENT, lit: ""},
		{in: `/*nice another\n try*/`, tok: cypher.COMMENT, lit: ""},
		{in: `/`, tok: cypher.DIV, lit: ""},
		{in: `  `, tok: cypher.WS, lit: "  "},
		{in: `[`, tok: cypher.LBRACKET, lit: ""},
	} {
		s := cypher.NewScanner(strings.NewReader(tc.in))
		tok, _, lit := s.Scan()
		if tok != tc.tok {
			t.Errorf("For input `%s` expected token '%s' got '%s' (%s)", tc.in, tc.tok, tok, lit)
		} else if lit != tc.lit {
			t.Errorf("Expected literal '%s' got '%s'", tc.in, lit)
		}
	}
}

func TestScanMultiple(t *testing.T) {
	type result struct {
		tok cypher.Token
		pos cypher.Pos
		lit string
	}
	exp := []result{
		{tok: cypher.MATCH, pos: cypher.Pos{Line: 0, Char: 0}, lit: ""},
		{tok: cypher.WS, pos: cypher.Pos{Line: 0, Char: 5}, lit: " "},
		{tok: cypher.LPAREN, pos: cypher.Pos{Line: 0, Char: 6}, lit: ""},
		{tok: cypher.IDENT, pos: cypher.Pos{Line: 0, Char: 7}, lit: "n"},
		{tok: cypher.COLON, pos: cypher.Pos{Line: 0, Char: 8}, lit: ""},
		{tok: cypher.IDENT, pos: cypher.Pos{Line: 0, Char: 9}, lit: "Person"},
		{tok: cypher.RPAREN, pos: cypher.Pos{Line: 0, Char: 15}, lit: ""},
		{tok: cypher.WS, pos: cypher.Pos{Line: 0, Char: 16}, lit: " "},
		{tok: cypher.WHERE, pos: cypher.Pos{Line: 0, Char: 17}, lit: ""},
		{tok: cypher.WS, pos: cypher.Pos{Line: 0, Char: 22}, lit: " "},
		{tok: cypher.IDENT, pos: cypher.Pos{Line: 0, Char: 23}, lit: "n"},
		{tok: cypher.DOT, pos: cypher.Pos{Line: 0, Char: 24}, lit: ""},
		{tok: cypher.IDENT, pos: cypher.Pos{Line: 0, Char: 25}, lit: "name"},
		{tok: cypher.WS, pos: cypher.Pos{Line: 0, Char: 29}, lit: " "},
		{tok: cypher.EQ, pos: cypher.Pos{Line: 0, Char: 30}, lit: ""},
		{tok: cypher.WS, pos: cypher.Pos{Line: 0, Char: 31}, lit: " "},
		{tok: cypher.STRING, pos: cypher.Pos{Line: 0, Char: 31}, lit: "Rafael"},
		{tok: cypher.WS, pos: cypher.Pos{Line: 0, Char: 40}, lit: " "},
		{tok: cypher.RETURN, pos: cypher.Pos{Line: 0, Char: 41}, lit: ""},
		{tok: cypher.WS, pos: cypher.Pos{Line: 0, Char: 47}, lit: " "},
		{tok: cypher.IDENT, pos: cypher.Pos{Line: 0, Char: 48}, lit: "n"},
		{tok: cypher.EOF, pos: cypher.Pos{Line: 0, Char: 50}, lit: ""},
	}

	// Create a scanner.
	v := `MATCH (n:Person) WHERE n.name = "Rafael" RETURN n`
	s := cypher.NewScanner(strings.NewReader(v))

	// Continually scan until we reach the end.
	var act []result
	for {
		tok, pos, lit := s.Scan()
		act = append(act, result{tok, pos, lit})
		if tok == cypher.EOF {
			break
		}
	}

	// Verify the token counts match.
	if len(exp) != len(act) {
		t.Fatalf("token count mismatch: exp=%d, got=%d", len(exp), len(act))
	}

	// Verify each token matches.
	for i := range exp {
		if !reflect.DeepEqual(exp[i], act[i]) {
			t.Fatalf("%d. token mismatch:\n\nexp=%#v\n\ngot=%#v", i, exp[i], act[i])
		}
	}
}
