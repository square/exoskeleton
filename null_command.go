package exoskeleton

import (
	"github.com/square/exit"
	"github.com/square/exoskeleton/v2/pkg/shellcomp"
)

// nullCommand represents an unrecognized Command.
type nullCommand struct {
	parent Command
	name   string
}

func (n nullCommand) Parent() Command          { return n.parent }
func (_ nullCommand) Path() string             { return "" }
func (n nullCommand) Name() string             { return n.name }
func (_ nullCommand) Summary() (string, error) { panic("Unused") }
func (_ nullCommand) Help() (string, error)    { panic("Unused") }

func (n nullCommand) Exec(e *Entrypoint, _, _ []string) error {
	return exit.ErrUnknownSubcommand
}

func (_ nullCommand) Complete(_ *Entrypoint, _, _ []string) ([]string, shellcomp.Directive, error) {
	// Unable to find the real command. E.g., <program> someInvalidCmd <TAB>
	return nil, shellcomp.DirectiveNoFileComp, nil
}

func (_ nullCommand) Subcommands() (Commands, error) {
	return Commands{}, nil
}
