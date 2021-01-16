package cli

import "flag"

type SubCommand struct {
	Arguments    []Argument
	ErrorHandler flag.ErrorHandling
}
