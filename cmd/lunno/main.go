package main

import (
	"lunno/cmd/lunno/cli"
	"os"
)

func main() {
	cli.Execute(os.Args[1:])
}
