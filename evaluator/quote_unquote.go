package evaluator

import (
	"fmt"

	"github.com/rtfb/tarsier/ast"
	"github.com/rtfb/tarsier/object"
	"github.com/rtfb/tarsier/token"
)

func quote(node ast.Node, env *object.Env) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{
		Node: node,
	}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Env) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}
		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}
		if len(call.Arguments) != 1 {
			return node
		}
		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToASTNode(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExpression.Function.TokenLiteral() == "unquote"
}

// TODO: here we create new tokens on the fly. That’s not a problem at the
// moment, but when our tokens will contain information about their origin,
// such as filename or line number, then we’ll also have to update these here,
// which might be quite difficult for tokens that are created dynamically.
func convertObjectToASTNode(o object.Object) ast.Node {
	switch o := o.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.Num,
			Literal: fmt.Sprintf("%d", o.Value),
		}
		return &ast.IntegerLiteral{
			Token: t,
			Value: o.Value,
		}
	case *object.Boolean:
		var t token.Token
		if o.Value {
			t = token.Token{
				Type:    token.True,
				Literal: "true",
			}
		} else {
			t = token.Token{
				Type:    token.False,
				Literal: "false",
			}
		}
		return &ast.Boolean{
			Token: t,
			Value: o.Value,
		}
	case *object.Quote:
		return o.Node
	// TODO: add the remaining types here
	default:
		return nil
	}
}
