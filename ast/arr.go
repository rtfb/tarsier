package ast

import (
	"bytes"
	"strings"

	"github.com/rtfb/tarsier/token"
)

// ArrayLiteral is the AST subtree containing an array literal.
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral implements Node.
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

// String implements Node.
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := make([]string, len(al.Elements))
	for i, el := range al.Elements {
		elements[i] = el.String()
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
