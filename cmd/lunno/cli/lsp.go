package cli

import (
	"flag"
	"lunno/internal/lsp"
)

type LspCommand struct{}

func (c *LspCommand) Name() string {
	return "lsp"
}

func (c *LspCommand) Description() string {
	return "Start lsp server."
}

func (c *LspCommand) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet(c.Name(), flag.ExitOnError)
}

func (c *LspCommand) Run(args []string) {
	lsp.StartLsp()
}
