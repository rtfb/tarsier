package ast

import (
	"github.com/rtfb/tarsier/token"
)

// IntegerLiteral is the AST subtree containing an integer literal.
type IntegerLiteral struct {
	Token token.Token // the first token of the expression
	Value int64
}

func (il *IntegerLiteral) expressionNode() {
}

// TokenLiteral implements Node.
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// String implements Node.
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
