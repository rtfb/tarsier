package ast

import (
	"bytes"
)

// Program is the root of the AST.
type Program struct {
	Statements []Statement
}

// TokenLiteral implements Node.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

// String implements Node.
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
