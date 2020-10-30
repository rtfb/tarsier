package evaluator

import (
	"github.com/rtfb/tarsier/ast"
	"github.com/rtfb/tarsier/object"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{
		Node: node,
	}
}
