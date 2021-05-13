package main

import (
	"fmt"
	"github.com/justclimber/fda-lang/fdalang"
	"io/ioutil"
	"log"
)

func main() {
	sourceCode, _ := ioutil.ReadFile("example/example1")
	fmt.Printf("Running source code:\n%s\n", string(sourceCode))
	l := fdalang.NewLexer(string(sourceCode))
	p, err := fdalang.NewParser(l)
	if err != nil {
		log.Fatalf("Lexing error: %s\n", err.Error())
	}

	astProgram, err := p.Parse()
	if err != nil {
		log.Fatalf("Parsing error: %s\n", err.Error())
	}
	env := fdalang.NewEnvironment()
	fmt.Println("Program output:")
	err = fdalang.NewExecAstVisitor().ExecAst(astProgram, env)
	if err != nil {
		log.Fatalf("Runtime error: %s\n", err.Error())
	}
	env.Print()
}
