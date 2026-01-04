package cli

import (
	"flag"
)

type Command interface {
	Name() string
	Description() string
	Run(args []string)
	FlagSet() *flag.FlagSet
}
