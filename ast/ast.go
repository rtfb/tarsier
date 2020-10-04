package ast

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
