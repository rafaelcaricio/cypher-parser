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

	ADD:        "ADD",
	ALL:        "ALL",
	AS:         "AS",
	ASC:        "ASC",
	ASCENDING:  "ASCENDING",
	BY:         "BY",
	CASE:       "CASE",
	CONSTRAINT: "CONSTRAINT",
	CONTAINS:   "CONTAINS",
	CREATE:     "CREATE",
	DELETE:     "DELETE",
	DESC:       "DESC",
	DESCENDING: "DESCENDING",
	DETACH:     "DETACH",
	DISTINCT:   "DISTINCT",
	DO:         "DO",
	DROP:       "DROP",
	ELSE:       "ELSE",
	END:        "END",
	ENDS:       "ENDS",
	EXISTS:     "EXISTS",
	FOR:        "FOR",
	IN:         "IN",
	IS:         "IS",
	LIMIT:      "LIMIT",
	MANDATORY:  "MANDATORY",
	MATCH:      "MATCH",
	MERGE:      "MERGE",
	OF:         "OF",
	ON:         "ON",
	OPTIONAL:   "OPTIONAL",
	ORDER:      "ORDER",
	REMOVE:     "REMOVE",
	REQUIRE:    "REQUIRE",
	RETURN:     "RETURN",
	SCALAR:     "SCALAR",
	SET:        "SET",
	SKIP:       "SKIP",
	STARTS:     "STARTS",
	THEN:       "THEN",
	UNION:      "UNION",
	UNIQUE:     "UNIQUE",
	UNWIND:     "UNWIND",
	WHEN:       "WHEN",
	WHERE:      "WHERE",
	WITH:       "WITH",
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

// tokstr returns a literal if provided, otherwise returns the token string.
func tokstr(tok Token, lit string) string {
	if lit != "" {
		return lit
	}
	return tok.String()
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
