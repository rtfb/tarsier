package lexer

import (
	"testing"

	"github.com/rtfb/tarsier/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
macro(x, y) { x + y; };
`
	tests := []struct {
		wantType    token.Type
		wantLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.Num, "5"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.Num, "10"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.LParen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.RParen, ")"},
		{token.Semicolon, ";"},
		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Asterisk, "*"},
		{token.Num, "5"},
		{token.Semicolon, ";"},
		{token.Num, "5"},
		{token.LT, "<"},
		{token.Num, "10"},
		{token.GT, ">"},
		{token.Num, "5"},
		{token.Semicolon, ";"},
		{token.If, "if"},
		{token.LParen, "("},
		{token.Num, "5"},
		{token.LT, "<"},
		{token.Num, "10"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Else, "else"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Num, "10"},
		{token.Equals, "=="},
		{token.Num, "10"},
		{token.Semicolon, ";"},
		{token.Num, "10"},
		{token.NotEquals, "!="},
		{token.Num, "9"},
		{token.Semicolon, ";"},
		{token.String, "foobar"},
		{token.String, "foo bar"},
		{token.LBracket, "["},
		{token.Num, "1"},
		{token.Comma, ","},
		{token.Num, "2"},
		{token.RBracket, "]"},
		{token.Semicolon, ";"},
		{token.LBrace, "{"},
		{token.String, "foo"},
		{token.Colon, ":"},
		{token.String, "bar"},
		{token.RBrace, "}"},
		{token.Macro, "macro"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.wantType {
			t.Fatalf("tests[%d] - wrong token type. want=%q, got=%q",
				i, tt.wantType, tok.Type)
		}
		if tok.Literal != tt.wantLiteral {
			t.Fatalf("tests[%d] - wrong literal. want=%q, got=%q",
				i, tt.wantLiteral, tok.Literal)
		}
	}
}
