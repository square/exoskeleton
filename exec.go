package exoskeleton

import (
	"os"

	"github.com/square/exit"
)

// Exec constructs an Entrypoint with the given paths and options and executes it.
func Exec(paths []string, options ...Option) {
	// Create a new Commandline application that will look for subcommands in the given paths.
	e, err := New(paths, options...)
	if err != nil {
		panic(err)
	}

	// Identify the subcommand being invoked from the arguments.
	cmd, args := e.Identify(os.Args[1:])

	// Execute the subcommand.
	err = cmd.Exec(e, args, os.Environ())

	// Exit the program with the exit code the subcommand returned.
	os.Exit(exit.FromError(err))
}
