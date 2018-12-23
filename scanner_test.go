package cypher

import (
	"reflect"
	"strings"
	"testing"
)

func TestScanInput(t *testing.T) {
	for _, tc := range []struct {
		in  string
		tok Token
		lit string
	}{
		{in: `something`, tok: IDENT, lit: "something"},
		{in: `desc`, tok: DESC, lit: ""},
		{in: `match`, tok: MATCH, lit: ""},
		{in: `or`, tok: OR, lit: ""},
		{in: `1233`, tok: INTEGER, lit: "1233"},
		{in: `3.14`, tok: NUMBER, lit: "3.14"},
		{in: `true`, tok: TRUE, lit: ""},
		{in: `null`, tok: NULL, lit: ""},
		{in: `"Hello, world!"`, tok: STRING, lit: "Hello, world!"},
		{in: `"String\nwith\nnewline"`, tok: STRING, lit: "String\nwith\nnewline"},
		{in: `'String\n'`, tok: STRING, lit: "String\n"},
		{in: ``, tok: EOF, lit: ""},
		{in: `<>`, tok: NEQ, lit: ""},
		{in: `<`, tok: LT, lit: ""},
		{in: `<=1`, tok: LTE, lit: ""},
		{in: `>`, tok: GT, lit: ""},
		{in: `>=`, tok: GTE, lit: ""},
		{in: `..`, tok: DOUBLEDOT, lit: ""},
		{in: `+`, tok: PLUS, lit: ""},
		{in: `+=`, tok: INC, lit: ""},
		{in: `//nice try`, tok: COMMENT, lit: ""},
		{in: `/*nice another\n try*/`, tok: COMMENT, lit: ""},
		{in: `/`, tok: DIV, lit: ""},
		{in: `  `, tok: WS, lit: "  "},
		{in: `[`, tok: LBRACKET, lit: ""},
	} {
		s := NewScanner(strings.NewReader(tc.in))
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
		tok Token
		pos Pos
		lit string
	}
	exp := []result{
		{tok: MATCH, pos: Pos{Line: 0, Char: 0}, lit: ""},
		{tok: WS, pos: Pos{Line: 0, Char: 5}, lit: " "},
		{tok: LPAREN, pos: Pos{Line: 0, Char: 6}, lit: ""},
		{tok: IDENT, pos: Pos{Line: 0, Char: 7}, lit: "n"},
		{tok: COLON, pos: Pos{Line: 0, Char: 8}, lit: ""},
		{tok: IDENT, pos: Pos{Line: 0, Char: 9}, lit: "Person"},
		{tok: RPAREN, pos: Pos{Line: 0, Char: 15}, lit: ""},
		{tok: WS, pos: Pos{Line: 0, Char: 16}, lit: " "},
		{tok: WHERE, pos: Pos{Line: 0, Char: 17}, lit: ""},
		{tok: WS, pos: Pos{Line: 0, Char: 22}, lit: " "},
		{tok: IDENT, pos: Pos{Line: 0, Char: 23}, lit: "n"},
		{tok: DOT, pos: Pos{Line: 0, Char: 24}, lit: ""},
		{tok: IDENT, pos: Pos{Line: 0, Char: 25}, lit: "name"},
		{tok: WS, pos: Pos{Line: 0, Char: 29}, lit: " "},
		{tok: EQ, pos: Pos{Line: 0, Char: 30}, lit: ""},
		{tok: WS, pos: Pos{Line: 0, Char: 31}, lit: " "},
		{tok: STRING, pos: Pos{Line: 0, Char: 31}, lit: "Rafael"},
		{tok: WS, pos: Pos{Line: 0, Char: 40}, lit: " "},
		{tok: RETURN, pos: Pos{Line: 0, Char: 41}, lit: ""},
		{tok: WS, pos: Pos{Line: 0, Char: 47}, lit: " "},
		{tok: IDENT, pos: Pos{Line: 0, Char: 48}, lit: "n"},
		{tok: EOF, pos: Pos{Line: 0, Char: 50}, lit: ""},
	}

	// Create a scanner.
	v := `MATCH (n:Person) WHERE n.name = "Rafael" RETURN n`
	s := NewScanner(strings.NewReader(v))

	// Continually scan until we reach the end.
	var act []result
	for {
		tok, pos, lit := s.Scan()
		act = append(act, result{tok, pos, lit})
		if tok == EOF {
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
