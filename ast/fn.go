package ast

import (
	"bytes"
	"strings"

	"github.com/rtfb/tarsier/token"
)

// FunctionLiteral represents the fn expression.
type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {
}

// TokenLiteral implements Node.
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// String implements Node.
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := make([]string, len(fl.Parameters))
	for i, p := range fl.Parameters {
		params[i] = p.String()
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}
