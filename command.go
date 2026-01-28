package exoskeleton

import (
	"github.com/square/exoskeleton/v2/pkg/shellcomp"
)

// Command is a command in your CLI application.
//
// In the Go CLI, 'go test' and 'go mod tidy' are both commands. Exoskeleton would
// model 'test' and 'tidy' as Commands. 'tidy' would be nested one level deeper
// than 'test', beneath a Command named 'mod'.
type Command interface {
	// Path returns the location of the executable that defines the command.
	// For built-in commands, it returns the path to the entrypoint executable itself.
	// It is used by the built-in command 'which'.
	Path() string

	// Name returns the name of the command.
	Name() string

	// Parent returns the command that contains this command.
	//
	// For unnested commands, Parent returns the Entrypoint. For the Entrypoint,
	// Parent returns nil.
	//
	// In the Go CLI, 'go test' and 'go mod tidy' are both commands. If modeled with
	// Exoskeleton, 'tidy''s Parent would be the 'mod' command and 'test''s Parent
	// would be the entrypoint itself, 'go'.
	Parent() Command

	// Exec executes the command.
	Exec(e *Entrypoint, args, env []string) error

	// Complete asks the command to return completions.
	// It is used by the built-in command 'complete' which provides shell completions.
	Complete(e *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error)

	// Summary returns the (short!) description of the command to be displayed
	// in menus.
	//
	// Returns a CommandError if the command does not fulfill the contract
	// for providing its summary.
	Summary() (string, error)

	// Help returns the help text for the command.
	//
	// Returns a CommandError if the command does not fulfill the contract
	// for providing its help.
	Help() (string, error)

	// Subcommands returns the list of Commands contained by this command.
	// Returns an empty slice for leaf commands (commands without subcommands).
	//
	// Returns a CommandError if the command does not fulfill the contract
	// for providing its subcommands.
	Subcommands() (Commands, error)
}
