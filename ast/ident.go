package ast

import (
	"github.com/rtfb/tarsier/token"
)

// Identifier is the AST subtree containing an identifier.
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {
}

// TokenLiteral implements Node.
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String implements Node.
func (i *Identifier) String() string {
	return i.Value
}
