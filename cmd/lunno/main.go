package main

import (
	"fmt"
	"lunno/internal/lexer"
	"lunno/internal/parser"
	"os"
)

func loadAndParse(filename string) *parser.Program {
	source, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	tokens, err := lexer.Tokenize(string(source), filename)
	if err != nil {
		return nil
	}

	program, errs := parser.ParseProgram(tokens)
	if len(errs) > 0 {
		fmt.Println("Parse errors:")
		for _, e := range errs {
			fmt.Println("  ", e)
		}
		return nil
	}
	fmt.Println(parser.DumpProgram(program))
	return program
}

func main() {
	loadAndParse(os.Args[1])
}
