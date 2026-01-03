package main

import (
	"fmt"
	"lunno/internal/lexer"
	"os"
)

func loadAndParse(filename string) {
	source, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	tokens, err := lexer.Tokenize(string(source), filename)
	if err != nil {
		return
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func main() {
	loadAndParse(os.Args[1])
}
