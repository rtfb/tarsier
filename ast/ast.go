package ast

import "github.com/rtfb/tarsier/token"

// Node represents any node of the AST.
type Node interface {
	TokenLiteral() string
}

// Statement is the statement kind of the AST node: statements do not produce a
// value.
type Statement interface {
	Node
	statementNode()
}

// Expression is the expression king of he AST node: expressions produce values.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root of the AST.
type Program struct {
	Statements []Statement
}

// TokenLiteral implements Node.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

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

// Identifier is the AST subtree containing an identifier.
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {
}

// TokenLiteral implements Node.
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
