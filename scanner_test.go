package cypher

import (
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
		{in: ``, tok: EOF, lit: ""},
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
