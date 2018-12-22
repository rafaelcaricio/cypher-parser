package cypher

import "testing"

func TestLookupToken(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected Token
	}{
		{"abc", IDENT},
		{"null", NULL},
		{"do", DO},
		{"not", NOT},
		{"unique", UNIQUE},
		{"starts", STARTS},
	} {
		if v := Lookup(tc.input); v != tc.expected {
			t.Errorf("Expected token '%s' got '%s'", tokens[tc.expected], tokens[v])
		}
	}
}
