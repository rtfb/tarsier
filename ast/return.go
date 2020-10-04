package ast

import (
	"bytes"

	"github.com/rtfb/tarsier/token"
)

// ReturnStatement is the AST subtree containing a return statement.
type ReturnStatement struct {
	Token       token.Token // the 'return' token itself
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {
}

// TokenLiteral implements Node.
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String implements Node.
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}
