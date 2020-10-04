package ast

import (
	"bytes"

	"github.com/rtfb/tarsier/token"
)

// PrefixExpression is the AST subtree containing a prefix expression.
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {
}

// TokenLiteral implements Node.
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String implements Node.
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}
