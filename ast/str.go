package ast

import (
	"github.com/rtfb/tarsier/token"
)

// StringLiteral represents a literal string.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral implements Node.
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

// String implements Node.
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}
