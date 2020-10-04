package ast

import (
	"bytes"

	"github.com/rtfb/tarsier/token"
)

// Node represents any node of the AST.
type Node interface {
	TokenLiteral() string
	String() string
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

// String implements Node.
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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

// String implements Node.
func (i *Identifier) String() string {
	return i.Value
}
