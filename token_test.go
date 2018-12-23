package cypher_test

import (
	"testing"

	"github.com/rafaelcaricio/cypher-parser"
)

func TestLookupToken(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected cypher.Token
	}{
		{"abc", cypher.IDENT},
		{"null", cypher.NULL},
		{"do", cypher.DO},
		{"not", cypher.NOT},
		{"unique", cypher.UNIQUE},
		{"starts", cypher.STARTS},
	} {
		if v := cypher.Lookup(tc.input); v != tc.expected {
			t.Errorf("Expected token '%s' got '%s'", tc.expected, v)
		}
	}
}
