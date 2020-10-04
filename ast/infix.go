package ast

import (
	"bytes"

	"github.com/rtfb/tarsier/token"
)

// InfixExpression is the AST subtree containing a infix expression.
type InfixExpression struct {
	Token    token.Token // the operator token, e.g. '+'
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {
}

// TokenLiteral implements Node.
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String implements Node.
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}
