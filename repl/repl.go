package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/rtfb/tarsier/evaluator"
	"github.com/rtfb/tarsier/lexer"
	"github.com/rtfb/tarsier/object"
	"github.com/rtfb/tarsier/parser"
)

const Prompt = ">> "

var stdlibFiles = []string{
	"stdlib/arr.ts",
	"stdlib/unless.ts",
}

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
	if err := doStdlib(stdlibFiles, out, env, macroEnv); err != nil {
		return err
	}
	return doFile(in, out, env, macroEnv)
}

func doStdlib(files []string, out io.Writer, env, macroEnv *object.Env) error {
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		if err := doFile(f, out, env, macroEnv); err != nil {
			return err
		}
	}
	return nil
}

func doFile(in io.Reader, out io.Writer, env, macroEnv *object.Env) error {
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
