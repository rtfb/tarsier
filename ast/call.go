package ast

import (
	"bytes"
	"strings"

	"github.com/rtfb/tarsier/token"
)

// CallExpression represents the call of a function.
type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {
}

// TokenLiteral implements Node.
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// String implements Node.
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := make([]string, len(ce.Arguments))
	for i, a := range ce.Arguments {
		args[i] = a.String()
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
