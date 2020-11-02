package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/rtfb/tarsier/evaluator"
	"github.com/rtfb/tarsier/lexer"
	"github.com/rtfb/tarsier/object"
	"github.com/rtfb/tarsier/parser"
)

const prompt = ">> "

// Start starts an interactive REPL.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv()
	macroEnv := object.NewEnv()
	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluator.DefineMacros(program, macroEnv)
		expandedProgram := evaluator.ExpandMacros(program, macroEnv)
		evaluated := evaluator.Eval(expandedProgram, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
