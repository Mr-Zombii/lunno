package cli

import (
	"fmt"
	"os"
)

func Execute(args []string) {
	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "-h", "--help", "help":
		printHelp()
		return
	case "-v", "--version", "version":
		findCommand("version").Run(nil)
		return
	}
	cmd := findCommand(args[0])
	if cmd == nil {
		_, err := fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
		if err != nil {
			return
		}
		printHelp()
		os.Exit(1)
	}

	cmd.Run(args[1:])
}
