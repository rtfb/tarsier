package ast

import (
	"bytes"
	"strings"

	"github.com/rtfb/tarsier/token"
)

// HashLiteral is the AST subtree containing a hash map literal.
type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

// TokenLiteral implements Node.
func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

// String implements Node.
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := make([]string, 0, len(hl.Pairs))
	for k, v := range hl.Pairs {
		pairs = append(pairs, k.String()+":"+v.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
