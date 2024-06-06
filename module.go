package exoskeleton

// Module is a namespace of commands in your CLI application.
//
// In the Go CLI, 'go mod tidy' is a command. Exoskeleton would module 'mod'
// as a Module and 'tidy' as a Command that's nested beneath it.
//
// Executing a Module produces a menu of its subcommands.
type Module interface {
	Command

	// Subcommands returns the list of Commands that are contained by this module.
	//
	// For example, in the Go CLI, 'go mod' is a Module and its Subcommands would
	// be 'download', 'edit', 'graph', 'init', 'tidy', 'vendor', 'verify', and 'why'.
	//
	// Returns a CommandError if the command does not fulfill the contract
	// for providing its subcommands.
	Subcommands() (Commands, error)
}
