package ast

import (
	"bytes"

	"github.com/rtfb/tarsier/token"
)

// BlockStatement represents any block statement.
type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {
}

// TokenLiteral implements Node.
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String implements Node.
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
