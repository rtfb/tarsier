package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/rtfb/tarsier/evaluator"
	"github.com/rtfb/tarsier/lexer"
	"github.com/rtfb/tarsier/object"
	"github.com/rtfb/tarsier/parser"
)

const Prompt = ">> "

// Start starts an interactive REPL.
func Start(in io.Reader, out io.Writer, prompt string) {
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

// DoFile interprets a program from a given Reader.
func DoFile(in io.Reader, out io.Writer) error {
	env := object.NewEnv()
	macroEnv := object.NewEnv()
	input, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	l := lexer.New(string(input))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return errors.New("TODO")
	}
	evaluator.DefineMacros(program, macroEnv)
	expandedProgram := evaluator.ExpandMacros(program, macroEnv)
	evaluated := evaluator.Eval(expandedProgram, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
	return nil
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
