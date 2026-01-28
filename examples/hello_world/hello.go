package main

import (
	"os"

	"github.com/square/exit"
	"github.com/square/exoskeleton/v2"
)

func main() {
	path, _ := os.Getwd()
	paths := []string{path + "/libexec"}

	// Create a new Commandline application that will look for subcommands
	// in the path './libexec'.
	cli, err := exoskeleton.New(paths)
	if err != nil {
		panic(err)
	}

	// Identify the subcommand being invoked from the arguments.
	cmd, args, err := cli.Identify(os.Args[1:])
	if err != nil {
		panic(err)
	}

	// Execute the subcommand.
	err = cmd.Exec(cli, args, os.Environ())

	// Exit the program with the exit code the subcommand returned.
	os.Exit(exit.FromError(err))
}
