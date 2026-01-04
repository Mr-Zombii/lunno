package cli

import (
	"fmt"
)

var commands = []Command{
	&RunCommand{},
	&VersionCommand{},
}

func findCommand(name string) Command {
	for _, cmd := range commands {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

func printHelp() {
	fmt.Println("Lunno - A small functional language")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  lunno <command> [flags] [args]")
	fmt.Println()
	fmt.Println("Available commands:")
	for _, cmd := range commands {
		fmt.Printf("  %-10s %s\n", cmd.Name(), cmd.Description())
	}
	fmt.Println()
	fmt.Println("Use 'lunno <command> -h' for command-specific help")
}
