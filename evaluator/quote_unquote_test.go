package evaluator

import (
	"testing"

	"github.com/rtfb/tarsier/object"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			"quote(5)",
			"5",
		},
		{
			"quote(5 + 8)",
			"(5 + 8)",
		},
		{
			"quote(foobar)",
			"foobar",
		},
		{
			"quote(foobar + barfoo)",
			"(foobar + barfoo)",
		},
		// XXX: this (invalid) test case causes the evaluator to panic. Need to
		// make it more robust.
		// {
		// 	"quote(foobar + barfoo",
		// 	"(foobar + barfoo)",
		// },
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote, got=%T (%+v)", evaluated, evaluated)
		}
		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}
		if quote.Node.String() != tt.want {
			t.Errorf("not equal, got=%q, want=%q", quote.Node.String(), tt.want)
		}
	}
}
