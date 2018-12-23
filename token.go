package cypher

import "strings"

// Token is a lexical token of the Cypher language.
type Token int

const (
	// ILLEGAL Token, EOF, WS are Special Cypher tokens.
	ILLEGAL Token = iota
	EOF
	WS
	COMMENT

	literalBeg
	// IDENT and the following are literal tokens.
	IDENT     // main
	NUMBER    // 12345.67
	INTEGER   // 12345
	STRING    // "abc"
	BADSTRING // "abc
	BADESCAPE // "\q
	TRUE      // true
	FALSE     // false
	NULL      // null
	literalEnd

	operatorBeg
	PLUS // +
	SUB  // -
	MUL  // *
	DIV  // /
	MOD  // %
	POW  // ^
	EQ   // =
	NEQ  // <>
	LT   // <
	LTE  // <=
	GT   // >
	GTE  // >=
	INC  // +=
	BAR  // |

	AND // AND
	OR  // OR
	XOR // XOR
	NOT // NOT
	operatorEnd

	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	COMMA     // ,
	COLON     // :
	SEMICOLON // ;
	DOT       // .
	DOUBLEDOT // ..

	keywordBeg
	// ALL and the following are Cypher Keywords
	ADD
	ALL
	AS
	ASC
	ASCENDING
	BY
	CASE
	CONSTRAINT
	CONTAINS
	CREATE
	DELETE
	DESC
	DESCENDING
	DETACH
	DISTINCT
	DO
	DROP
	ELSE
	END
	ENDS
	EXISTS
	FOR
	IN
	IS
	LIMIT
	MANDATORY
	MATCH
	MERGE
	OF
	ON
	OPTIONAL
	ORDER
	REMOVE
	REQUIRE
	RETURN
	SCALAR
	SET
	SKIP
	STARTS
	THEN
	UNION
	UNIQUE
	UNWIND
	WHEN
	WHERE
	WITH
	keywordEnd
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	WS:      "WS",

	IDENT:  "IDENT",
	NUMBER: "NUMBER",
	STRING: "STRING",
	TRUE:   "TRUE",
	FALSE:  "FALSE",

	PLUS: "+",
	SUB:  "-",
	MUL:  "*",
	DIV:  "/",
	MOD:  "%",
	POW:  "^",

	AND: "AND",
	OR:  "OR",
	XOR: "XOR",
	NOT: "NOT",

	EQ:  "=",
	NEQ: "<>",
	LT:  "<",
	LTE: "<=",
	GT:  ">",
	GTE: ">=",

	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACKET:  "[",
	RBRACKET:  "]",
	COMMA:     ",",
	COLON:     ":",
	SEMICOLON: ";",
	DOT:       ".",

	CREATE:     "CREATE",
	DELETE:     "DELETE",
	DETACH:     "DETACH",
	EXISTS:     "EXISTS",
	MATCH:      "MATCH",
	MERGE:      "MERGE",
	OPTIONAL:   "OPTIONAL",
	REMOVE:     "REMOVE",
	RETURN:     "RETURN",
	SET:        "SET",
	UNION:      "UNION",
	UNWIND:     "UNWIND",
	WITH:       "WITH",
	LIMIT:      "LIMIT",
	ORDER:      "ORDER",
	SKIP:       "SKIP",
	WHERE:      "WHERE",
	ASC:        "ASC",
	ASCENDING:  "ASCENDING",
	BY:         "BY",
	DESC:       "DESC",
	DESCENDING: "DESCENDING",
	ON:         "ON",
	ALL:        "ALL",
	CASE:       "CASE",
	ELSE:       "ELSE",
	END:        "END",
	THEN:       "THEN",
	WHEN:       "WHEN",
	AS:         "AS",
	CONTAINS:   "CONTAINS",
	DISTINCT:   "DISTINCT",
	ENDS:       "ENDS",
	IN:         "IN",
	IS:         "IS",
	STARTS:     "STARTS",
	CONSTRAINT: "CONSTRAINT",
	DO:         "DO",
	ADD:        "ADD",
	DROP:       "DROP",
	FOR:        "FOR",
	MANDATORY:  "MANDATORY",
	OF:         "OF",
	REQUIRE:    "REQUIRE",
	SCALAR:     "SCALAR",
	UNIQUE:     "UNIQUE",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for tok := keywordBeg + 1; tok < keywordEnd; tok++ {
		keywords[strings.ToLower(tokens[tok])] = tok
	}
	for _, tok := range []Token{AND, OR, XOR, NOT} {
		keywords[strings.ToLower(tokens[tok])] = tok
	}
	keywords["true"] = TRUE
	keywords["false"] = FALSE
	keywords["null"] = NULL
}

// isOperator returns true for operator tokens.
func (tok Token) isOperator() bool { return tok > operatorBeg && tok < operatorEnd }

// String returns the string representation of the token.
func (tok Token) String() string {
	if tok >= 0 && tok < Token(len(tokens)) {
		return tokens[tok]
	}
	return ""
}

// Lookup returns the token associated with a given string.
func Lookup(ident string) Token {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}

// Pos specifies the line and character position of a token.
// The Char and Line are both zero-based indexes.
type Pos struct {
	Line int
	Char int
}
