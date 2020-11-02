package evaluator

import (
	"testing"

	"github.com/rtfb/tarsier/ast"
	"github.com/rtfb/tarsier/lexer"
	"github.com/rtfb/tarsier/object"
	"github.com/rtfb/tarsier/parser"
)

func TestDefineMacros(t *testing.T) {
	input := `
	let number = 1;
	let function = fn(x, y) { x + y };
	let mymacro = macro(x, y) { x + y };
	`
	env := object.NewEnv()
	program := testParseProgram(input)
	DefineMacros(program, env)
	if len(program.Statements) != 2 {
		t.Fatalf("wrong number of statements, want=2, got=%d", len(program.Statements))
	}
	_, ok := env.Get("number")
	if ok {
		t.Fatalf("number should not be defined")
	}
	_, ok = env.Get("function")
	if ok {
		t.Fatalf("function should not be defined")
	}
	obj, ok := env.Get("mymacro")
	if !ok {
		t.Fatalf("macro not in invironment")
	}
	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("object is not Macro, got=%T (%+v)", obj, obj)
	}
	if len(macro.Parameters) != 2 {
		t.Fatalf("wrong number of macro parameters, want=2, got=%d",
			len(macro.Parameters))
	}
	if macro.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got=%q", macro.Parameters[0])
	}
	if macro.Parameters[1].String() != "y" {
		t.Fatalf("parameter is not 'y', got=%q", macro.Parameters[1])
	}
	wantBody := "(x + y)"
	if macro.Body.String() != wantBody {
		t.Fatalf("body is not %q, got=%q", wantBody, macro.Body.String())
	}
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			`let infixExpression = macro() { quote(1 + 2) };
			infixExpression();`,
			"(1 + 2)",
		},
		{
			`let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
			reverse(2 + 2, 10 - 5);`,
			"(10 - 5) - (2 + 2)",
		},
		{
			`let unless = macro(condition, consequence, alternative) {
				quote(if (!(unquote(condition))) {
					unquote(consequence);
				} else {
					unquote(alternative);
				});
			};
			unless(10 > 5, puts("not greater"), puts("greater"));`,
			`if (!(10 > 5)) { puts("not greater") } else { puts("greater") }`,
		},
	}
	for _, tt := range tests {
		want := testParseProgram(tt.want)
		program := testParseProgram(tt.input)
		env := object.NewEnv()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)
		if expanded.String() != want.String() {
			t.Errorf("not equal: want=%q, got=%q", want.String(), expanded.String())
		}
	}
}

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
