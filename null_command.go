package exoskeleton

import (
	"github.com/square/exit"
	"github.com/squareup/exoskeleton/pkg/shellcomp"
)

// nullCommand represents an unrecognized Command.
type nullCommand struct {
	parent Module
	name   string
}

func (n nullCommand) Parent() Module           { return n.parent }
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

// IsNull returns true if the given Command is a NullCommand and false if it is not.
func IsNull(command Command) bool {
	_, ok := command.(nullCommand)
	return ok
}
