package ast

import (
	"bytes"

	"github.com/rtfb/tarsier/token"
)

// LetStatement is the AST subtree containing a let statement.
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {
}

// TokenLiteral implements Node.
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// String implements Node.
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
