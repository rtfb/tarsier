package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/rtfb/tarsier/repl"
)

func main() {
	if len(os.Args) > 1 {
		runWithFileArg(os.Args[1])
		return
	}
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello, %s! This is the Tarsier programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout, repl.Prompt)
}

func runWithFileArg(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	return repl.DoFile(f, os.Stdout)
}
