package ast

import (
	"bytes"

	"github.com/rtfb/tarsier/token"
)

// IfExpression represents the if expression.
type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {
}

// TokenLiteral implements Node.
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String implements Node.
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}
