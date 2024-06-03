package exoskeleton

import (
	"github.com/square/exoskeleton/pkg/shellcomp"
)

// Command is a command in your CLI application.
//
// In the Go CLI, 'go test' and 'go mod tidy' are both commands. Exoskeleton would
// model 'test' and 'tidy' as Commands. 'tidy' would be nested one level deeper
// than 'test', beneath a Module named 'mod'.
type Command interface {
	// Path returns the location of the executable that defines the command.
	// For built-in commands, it returns the path to the entrypoint executable itself.
	// It is used by the built-in command 'which'.
	Path() string

	// Name returns the name of the command.
	Name() string

	// Parent returns the module that contains the command.
	//
	// For unnamespaced commands, Parent returns the Entrypoint. For the Entrypoint,
	// Parent returns nil.
	//
	// In the Go CLI, 'go test' and 'go mod tidy' are both commands. If modeled with
	// Exoskeleton, 'tidy''s Parent would be the 'mod' Module and 'test''s Parent
	// would be the entrypoint itself, 'go'.
	Parent() Module

	// Exec executes the command.
	Exec(e *Entrypoint, args, env []string) error

	// Complete asks the command to return completions.
	// It is used by the built-in command 'complete' which provides shell completions.
	Complete(e *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error)

	// Summary returns the (short!) description of the command to be displayed
	// in menus.
	Summary() string

	// Help returns the help text for the command.
	Help() string
}
