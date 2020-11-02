package evaluator

import (
	"github.com/rtfb/tarsier/ast"
	"github.com/rtfb/tarsier/object"
)

// DefineMacros extracts macro definitions from a given AST, and stores them in
// an environment.
func DefineMacros(program *ast.Program, env *object.Env) {
	definitions := []int{}
	for i, statement := range program.Statements {
		if isMacroDefinition(statement) {
			// NB: storing in the env happens inside this helper method:
			addMacro(statement, env)
			definitions = append(definitions, i)
		}
	}
	for i := len(definitions) - 1; i >= 0; i-- {
		definitionIndex := definitions[i]
		program.Statements = append(
			program.Statements[:definitionIndex],
			program.Statements[definitionIndex+1:]...,
		)
	}
}

func isMacroDefinition(node ast.Statement) bool {
	letStatement, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}
	_, ok = letStatement.Value.(*ast.MacroLiteral)
	return ok
}

func addMacro(stmt ast.Statement, env *object.Env) {
	letStmt, _ := stmt.(*ast.LetStatement)
	macroLiteral, _ := letStmt.Value.(*ast.MacroLiteral)
	macro := object.Macro{
		Parameters: macroLiteral.Parameters,
		Body:       macroLiteral.Body,
		Env:        env,
	}
	env.Set(letStmt.Name.Value, &macro)
}

// ExpandMacros takes an AST and an environment and expands macros found in the
// AST, reinserting the generated code back in.
func ExpandMacros(program ast.Node, env *object.Env) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}
		macro, ok := isMacroCall(callExpression, env)
		if !ok {
			return node
		}
		args := quoteArgs(callExpression)
		evalEnv := extendMacroEnv(macro, args)
		evaluated := Eval(macro.Body, evalEnv)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macros")
		}
		return quote.Node
	})
}

func isMacroCall(exp *ast.CallExpression, env *object.Env) (*object.Macro, bool) {
	identifier, ok := exp.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}
	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil, false
	}
	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}
	return macro, true
}

func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	args := make([]*object.Quote, len(exp.Arguments))
	for i, a := range exp.Arguments {
		args[i] = &object.Quote{Node: a}
	}
	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Env {
	extended := object.NewEnclosedEnv(macro.Env)
	for i, param := range macro.Parameters {
		extended.Set(param.Value, args[i])
	}
	return extended
}
