package ast

import (
	"github.com/rtfb/tarsier/token"
)

// Boolean is the AST subtree containing a boolean literal.
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {
}

// TokenLiteral implements Node.
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// String implements Node.
func (b *Boolean) String() string {
	return b.Token.Literal
}
