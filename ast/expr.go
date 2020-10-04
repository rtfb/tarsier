package ast

import (
	"github.com/rtfb/tarsier/token"
)

// ExpressionStatement is the AST subtree containing an expression.
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {
}

// TokenLiteral implements Node.
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String implements Node.
func (es *ExpressionStatement) String() string {
	if es.Expression == nil {
		return ""
	}
	return es.Expression.String()
}
