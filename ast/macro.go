package ast

import (
	"bytes"
	"strings"

	"github.com/rtfb/tarsier/token"
)

// MacroLiteral represents the definition of a macro.
type MacroLiteral struct {
	Token      token.Token // the 'macro' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (ml *MacroLiteral) expressionNode() {}

// TokenLiteral implements Node.
func (ml *MacroLiteral) TokenLiteral() string {
	return ml.Token.Literal
}

// String implements Node.
func (ml *MacroLiteral) String() string {
	var out bytes.Buffer
	params := make([]string, len(ml.Parameters))
	for i, p := range ml.Parameters {
		params[i] = p.String()
	}
	out.WriteString(ml.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString("( ")
	out.WriteString(ml.Body.String())
	return out.String()
}
