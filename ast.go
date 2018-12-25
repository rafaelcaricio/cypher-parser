package cypher

import (
	"bytes"
	"fmt"
)

// Query represents the Cypher query root element.
type Query struct {
	Root *SingleQuery
}

func (q Query) String() string {
	return q.Root.String()
}

// SingleQuery ...
type SingleQuery struct {
	Reading     []ReadingClause
	Distinct    bool
	ReturnItems []Expr
	Order       []OrderBy
	Skip        *Expr
	Limit       *Expr
}

func (sq SingleQuery) String() string {
	var buf bytes.Buffer

	for _, r := range sq.Reading {
		_, _ = buf.WriteString(r.String())
	}

	_, _ = buf.WriteString(" RETURN ")

	if sq.Distinct {
		_, _ = buf.WriteString("DISTINCT ")
	}

	for _, i := range sq.ReturnItems {
		_, _ = buf.WriteString(i.String())
	}

	if len(sq.Order) > 0 {
		_, _ = buf.WriteString(" ORDER BY ")
		for _, o := range sq.Order {
			_, _ = buf.WriteString(o.String())
		}
	}

	return buf.String()
}

// ReadingClause ...
type ReadingClause struct {
	OptionalMatch bool
	Pattern       []MatchPattern
	Where         *Expr
	// Unwind
	// Call
}

func (rc ReadingClause) String() string {
	var buf bytes.Buffer

	if rc.OptionalMatch {
		_, _ = buf.WriteString(" OPTIONAL ")
	}

	for _, p := range rc.Pattern {
		_, _ = buf.WriteString(p.String())
	}

	if w := rc.Where; w != nil {
		_, _ = buf.WriteString((*w).String())
	}

	return buf.String()
}

// MatchPattern ...
type MatchPattern struct {
	Variable *Variable
	Elements []PatternElement
}

func (mp MatchPattern) String() string {
	var buf bytes.Buffer

	buf.WriteString("MATCH ")

	if mp.Variable != nil {
		_, _ = buf.WriteString((*mp.Variable).String())
		_, _ = buf.WriteString(" = ")
	}

	for _, e := range mp.Elements {
		_, _ = buf.WriteString(e.String())
	}

	return buf.String()
}

// PatternElement ...
type PatternElement interface {
	patternElem()
	String() string
}

func (np NodePattern) patternElem() {}
func (ep EdgePattern) patternElem() {}

// NodePattern ...
type NodePattern struct {
	Variable   *Variable
	Labels     []string
	Properties map[string]Expr
}

func (np NodePattern) String() string {
	var buf bytes.Buffer

	_, _ = buf.WriteRune('(')

	if np.Variable != nil {
		_, _ = buf.WriteString((*np.Variable).String())
	}

	for _, l := range np.Labels {
		_, _ = buf.WriteString(" :")
		_, _ = buf.WriteString(l)
	}

	if len(np.Properties) > 0 {
		_, _ = buf.WriteString(" {")

		var next bool
		for p, v := range np.Properties {
			if next {
				_, _ = buf.WriteRune(',')
			}
			_, _ = buf.WriteString(p)
			_, _ = buf.WriteString(": ")
			_, _ = buf.WriteString(v.String())
			next = true
		}

		_, _ = buf.WriteRune('}')
	}

	_, _ = buf.WriteRune(')')

	return buf.String()
}

// EdgePattern ...
type EdgePattern struct {
	Variable   *string
	Labels     []string
	Properties map[string]Expr
	MinHops    *int
	MaxHops    *int
	Direction  EdgeDirection
}

// Var ...
func (ep EdgePattern) Var() *string {
	return ep.Variable
}

func (ep EdgePattern) String() string {
	var buf bytes.Buffer

	switch ep.Direction {
	case EdgeRight, EdgeUndefined:
		_, _ = buf.WriteRune('-')
	case EdgeLeft, EdgeOutgoing:
		_, _ = buf.WriteString("<-")
	}

	_, _ = buf.WriteRune('[')

	if ep.Variable != nil {
		_, _ = buf.WriteString(*ep.Variable)
	}

	for i, l := range ep.Labels {
		if i > 0 {
			_, _ = buf.WriteString(" | ")
		}
		_, _ = buf.WriteRune(':')
		_, _ = buf.WriteString(l)
	}

	if len(ep.Properties) > 0 {
		_, _ = buf.WriteRune('{')

		var next bool
		for p, v := range ep.Properties {
			if next {
				_, _ = buf.WriteRune(',')
			}
			_, _ = buf.WriteString(p)
			_, _ = buf.WriteRune(':')
			_, _ = buf.WriteString(v.String())
			next = true
		}

		_, _ = buf.WriteRune('}')
	}

	_, _ = buf.WriteRune(']')

	switch ep.Direction {
	case EdgeLeft, EdgeUndefined:
		_, _ = buf.WriteRune('-')
	case EdgeRight, EdgeOutgoing:
		_, _ = buf.WriteString("->")
	}

	return buf.String()
}

// EdgeDirection ...
type EdgeDirection int

const (
	EdgeUndefined EdgeDirection = iota
	EdgeRight
	EdgeLeft
	EdgeOutgoing
)

// OrderDirection ...
type OrderDirection int

const (
	// Ascending defines the ascending ordering.
	Ascending OrderDirection = iota
	// Descending defines the descending ordering.
	Descending
)

// OrderBy ...
type OrderBy struct {
	Dir  OrderDirection
	Item Expr
}

func (o OrderBy) String() string {
	return ""
}

// Variable ...
type Variable string

func (v Variable) String() string {
	return string(v)
}

// Symbol ...
type Symbol string

func (s Symbol) String() string {
	return string(s)
}

// StrLiteral ...
type StrLiteral string

func (s StrLiteral) String() string {
	return fmt.Sprintf("\"%s\"", string(s))
}

// Expr ...
type Expr interface {
	exp()
	String() string
}

func (v Variable) exp()   {}
func (s Symbol) exp()     {}
func (s StrLiteral) exp() {}
