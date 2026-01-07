package cli

import (
	"flag"
	"fmt"
	"lunno/internal/lexer"
	"os"
)

type RunCommand struct {
	dumpAST *bool
}

func (c *RunCommand) Name() string {
	return "run"
}

func (c *RunCommand) Description() string {
	return "Run a Lunno source file"
}

func (c *RunCommand) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.dumpAST = fs.Bool("dump-ast", false, "Print AST of program")
	return fs
}

func (c *RunCommand) Run(args []string) {
	fs := c.FlagSet()
	err := fs.Parse(args)
	if err != nil {
		return
	}
	files := fs.Args()
	if len(files) < 1 {
		fmt.Println("Please specify a source file to run")
		os.Exit(1)
	}
	filename := files[0]
	source, err := os.ReadFile(filename)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
	_, tokens, err := lexer.Tokenize(string(source), filename)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Lexing error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
	//program, errs := parser.ParseProgram(tokens, lx)
	//if len(errs) > 0 {
	//	fmt.Printf("Parse errors (%d):\n", len(errs))
	//	for _, e := range errs {
	//		fmt.Println(" ", e)
	//	}
	//	os.Exit(1)
	//}
	//typeErrors := typechecker.Check(program)
	//if len(typeErrors) > 0 {
	//	fmt.Printf("Type errors (%d):\n", len(typeErrors))
	//	for _, err := range typeErrors {
	//		fmt.Println(" ", err)
	//	}
	//	os.Exit(1)
	//}
	//if *c.dumpAST {
	//	fmt.Println(parser.DumpProgram(program))
	//	return
	//}
	for _, tok := range tokens {
		fmt.Println(tok)
	}
	fmt.Println("Program ran successfully!")
}
