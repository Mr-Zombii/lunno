package cli

import (
	"flag"
	"fmt"
	"lunno/internal/version"
)

type VersionCommand struct{}

func (c *VersionCommand) Name() string {
	return "version"
}

func (c *VersionCommand) Description() string {
	return "Show Lunno version"
}

func (c *VersionCommand) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet(c.Name(), flag.ExitOnError)
}

func (c *VersionCommand) Run(args []string) {
	fmt.Println(version.Full())
}
