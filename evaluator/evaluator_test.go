package evaluator

import (
	"testing"

	"github.com/rtfb/tarsier/lexer"
	"github.com/rtfb/tarsier/object"
	"github.com/rtfb/tarsier/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.want)
	}
}

func TestEvalBoolExpression(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBoolObject(t, evaluated, tt.want)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBoolObject(t, evaluated, tt.want)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		got := testEval(t, tt.input)
		integer, ok := tt.want.(int)
		if ok {
			testIntegerObject(t, got, int64(integer))
		} else {
			testNullObject(t, got)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return 10;
			}
			return 1;
		}`, 10,
		},
	}
	for _, tt := range tests {
		got := testEval(t, tt.input)
		testIntegerObject(t, got, tt.want)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input   string
		wantMsg string
	}{
		{
			"5 + true",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar;",
			`identifier not found: "foobar"`,
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
	}
	for _, tt := range tests {
		evaled := testEval(t, tt.input)
		errObj, ok := evaled.(*object.Error)
		if !ok {
			t.Errorf("no error object returned, got=%T (%+v)", evaled, evaled)
			continue
		}
		if errObj.Message != tt.wantMsg {
			t.Errorf("wrong error message, want=%q, got=%q", tt.wantMsg, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.want)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; }"
	evaluated := testEval(t, input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object it not Function, got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong number of params, got %+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got=%q", fn.Parameters[0])
	}
	wantBody := "(x + 2)"
	if fn.Body.String() != wantBody {
		t.Fatalf("body is not %q, got=%q", wantBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.want)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(2);`
	testIntegerObject(t, testEval(t, input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String, got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value, got=%q", str.Value)
	}
}

func TestStringConcat(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String, got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value, got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments, got=2, want=1"},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch want := tt.want.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(want))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object it not Error, git=%T (%+v)", evaluated, evaluated)
			}
			if errObj.Message != want {
				t.Errorf("wrong error message, want=%q, got=%q", want, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array, got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements, want=3, got=%d", len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i]",
			1,
		},
		{
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.want.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`
	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash, got=%T (%+v)", evaluated, evaluated)
	}
	want := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		True.HashKey():                             5,
		False.HashKey():                            6,
	}
	if len(result.Pairs) != len(want) {
		t.Fatalf("Hash has wrong num or pairs, got=%d, want=%d", len(result.Pairs),
			len(want))
	}
	for wantKey, wantValue := range want {
		pair, ok := result.Pairs[wantKey]
		if !ok {
			t.Errorf("no pair for key %q in Pairs", wantKey)
		}
		testIntegerObject(t, pair.Value, wantValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.want.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != Null {
		t.Errorf("object is not Null, got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	p.CheckParseErrors(t)
	env := object.NewEnv()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, want int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer, got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != want {
		t.Errorf("object has wrong value, want=%d, got=%d", want, result.Value)
		return false
	}
	return true
}

func testBoolObject(t *testing.T, obj object.Object, want bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean, got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != want {
		t.Errorf("object has wrong value, want=%t, got=%t", want, result.Value)
		return false
	}
	return true
}
