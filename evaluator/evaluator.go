package evaluator

import (
	"github.com/rtfb/tarsier/ast"
	"github.com/rtfb/tarsier/object"
)

// Eval evaluates an AST passed to it and returns an object it evaluates to.
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// statements:
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// expressions:
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}
