package ast

import (
	"bytes"

	"github.com/rtfb/tarsier/token"
)

// IndexExpression is the AST subtree containing an array indexing expression.
type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral implements Node.
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String implements Node.
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}
