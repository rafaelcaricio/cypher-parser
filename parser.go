package cypher

import (
	"fmt"
	"io"
	"strings"
)

// Parser represents a Cypher parser.
type Parser struct {
	s *bufScanner
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: newBufScanner(r)}
}

// ParseQuery parses a query string and returns its AST representation.
func ParseQuery(s string) (Query, error) {
	return NewParser(strings.NewReader(s)).ParseQuery()
}

// ParseQuery parses a Cypher string and returns a Query AST object.
func (p *Parser) ParseQuery() (q Query, err error) {
	for {
		if tok, _, _ := p.ScanIgnoreWhitespace(); tok == EOF {
			return q, nil
		} else if tok == SEMICOLON {
			continue
		} else {
			p.Unscan()
			sq, err := p.ParseSingleQuery()
			if err != nil {
				return q, err
			}
			q = Query{Root: *sq}
		}
	}
}

// ParseSingleQuery ...
func (p *Parser) ParseSingleQuery() (*SingleQuery, error) {
	// read all MATCH clauses
	sq := &SingleQuery{}
	for {
		if tok, _, _ := p.ScanIgnoreWhitespace(); tok != MATCH && tok != OPTIONAL {
			p.Unscan()
			break
		} else {
			p.Unscan()
			r, err := p.ScanReadingClause()
			if err != nil {
				return nil, err
			}
			sq.Reading = append(sq.Reading, *r)
		}
	}
	// scan return, if not it's an error because RETURN is obligatory
	if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != RETURN {
		return nil, newParseError(tokstr(tok, lit), []string{"RETURN"}, pos)
	}

	tok, _, _ := p.ScanIgnoreWhitespace()
	if tok == DISTINCT {
		sq.Distinct = true
	} else {
		p.Unscan()
	}

	return sq, nil
}

// ScanReadingClause ...
func (p *Parser) ScanReadingClause() (*ReadingClause, error) {
	rc := &ReadingClause{}

	// might be optionally matching this
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok == OPTIONAL {
		rc.OptionalMatch = true
	} else {
		p.Unscan()
	}

	// MATCH is obligatory here
	if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != MATCH {
		return nil, newParseError(tokstr(tok, lit), []string{"MATCH"}, pos)
	}

	for {
		mp, err := p.ScanMatchPattern()
		if err != nil {
			return nil, err
		}
		rc.Pattern = append(rc.Pattern, *mp)

		if tok, _, _ := p.ScanIgnoreWhitespace(); tok != COMMA {
			p.Unscan()
			break
		}
	}

	// might be optional WHERE
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok == WHERE {
		exp, err := p.ScanExpression()
		if err != nil {
			return nil, err
		}
		rc.Where = &exp
	} else {
		p.Unscan()
	}

	return rc, nil
}

// ScanMatchPattern ...
func (p *Parser) ScanMatchPattern() (*MatchPattern, error) {
	mp := &MatchPattern{}

	if tok, _, lit := p.ScanIgnoreWhitespace(); tok == IDENT {
		// We need the `=` character here
		if tok1, pos, lit1 := p.ScanIgnoreWhitespace(); tok1 != EQ {
			return nil, newParseError(tokstr(tok1, lit1), []string{"="}, pos)
		}

		v := Variable(lit)
		mp.Variable = &v
	} else {
		p.Unscan()
	}

	// scan the pattern itself
	elems, err := p.ScanPatternElements()
	if err != nil {
		return nil, err
	}
	mp.Elements = elems

	return mp, nil
}

// ScanPatternElements ...
func (p *Parser) ScanPatternElements() (pe []PatternElement, err error) {
	var node *NodePattern
	numParens := 0
	for {
		node, err = p.ScanNodePattern()
		if err != nil {
			return nil, err
		}
		if node == nil {
			// might be only parens around the actual match, lets try...
			if tok, pos, lit := p.ScanIgnoreWhitespace(); tok == LPAREN {
				numParens++
			} else {
				return nil, newParseError(tokstr(tok, lit), []string{"("}, pos)
			}
		} else {
			break
		}
	}

	pe = []PatternElement{node}

	for {
		edge, err := p.ScanEdgePattern()
		if err != nil {
			return nil, err
		} else if edge == nil {
			break
		} else {
			pe = append(pe, edge)
		}
		node, err = p.ScanNodePattern()
		if err != nil {
			return nil, err
		} else if node == nil {
			break
		} else {
			pe = append(pe, node)
		}
	}

	// need to close all open parens
	for i := 0; i < numParens; i++ {
		if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != RPAREN {
			return nil, newParseError(tokstr(tok, lit), []string{")"}, pos)
		}
	}

	return pe, nil
}

// ScanNodePattern returns a NodePattern if possible to consume a complete valid node.
func (p *Parser) ScanNodePattern() (*NodePattern, error) {
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok != LPAREN {
		// We already know we cannot consume a valid node if the pattern doesn't start with `(`
		p.Unscan()
		return nil, nil
	}
	var validNode bool
	var node NodePattern
	if tok, _, lit := p.ScanIgnoreWhitespace(); tok == IDENT {
		v := Variable(lit)
		node.Variable = &v
		validNode = true
	} else {
		p.Unscan()
	}

	for {
		if tok, _, _ := p.ScanIgnoreWhitespace(); tok == COLON {
			if tok1, pos, lit := p.ScanIgnoreWhitespace(); tok1 == IDENT {
				node.Labels = append(node.Labels, lit)
				validNode = true
			} else {
				return nil, newParseError(tokstr(tok, lit), []string{"Label Identifier"}, pos)
			}
		} else {
			p.Unscan()
			break
		}
	}

	props, err := p.ScanProperties()
	if err != nil {
		return nil, err
	} else if props != nil {
		node.Properties = *props
		validNode = true
	}

	if tok, pos, lit := p.ScanIgnoreWhitespace(); tok == RPAREN {
		return &node, nil
	} else if validNode && tok != RPAREN {
		// We need to close the node definition
		return nil, newParseError(tokstr(tok, lit), []string{")"}, pos)
	}

	p.Unscan()
	p.Unscan() // unscan the first LPAREN then
	return nil, nil
}

// ScanEdgePattern returns an EdgePattern if possible to consume a complete valid edge.
func (p *Parser) ScanEdgePattern() (*EdgePattern, error) {
	return nil, nil
}

// ScanExpression ...
func (p *Parser) ScanExpression() (Expr, error) {
	return Variable("p_test"), nil
}

// ScanProperties ...
func (p *Parser) ScanProperties() (*map[string]Expr, error) {
	return nil, nil
}

// Scan returns the next token from the underlying scanner.
func (p *Parser) Scan() (tok Token, pos Pos, lit string) { return p.s.Scan() }

// ScanIgnoreWhitespace scans the next non-whitespace and non-comment token.
func (p *Parser) ScanIgnoreWhitespace() (tok Token, pos Pos, lit string) {
	for {
		tok, pos, lit = p.Scan()
		if tok == WS || tok == COMMENT {
			continue
		}
		return
	}
}

// Unscan pushes the previously read token back onto the buffer.
func (p *Parser) Unscan() { p.s.Unscan() }

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Message  string
	Found    string
	Expected []string
	Pos      Pos
}

// newParseError returns a new instance of ParseError.
func newParseError(found string, expected []string, pos Pos) *ParseError {
	return &ParseError{Found: found, Expected: expected, Pos: pos}
}

// Error returns the string representation of the error.
func (e *ParseError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s at line %d, char %d", e.Message, e.Pos.Line+1, e.Pos.Char+1)
	}
	return fmt.Sprintf("found %s, expected %s at line %d, char %d", e.Found, strings.Join(e.Expected, ", "), e.Pos.Line+1, e.Pos.Char+1)
}
