package exoskeleton

// IsEmbedded returns true if the given Command is built into the exoskeleton.
func IsEmbedded(command Command) bool {
	_, ok := command.(*builtinCommand)
	return ok
}

// IsNull returns true if the given Command is a NullCommand and false if it is not.
func IsNull(command Command) bool {
	_, ok := command.(nullCommand)
	return ok
}

// SubcommandsReporter is implemented by Commands that can report whether
// they may have subcommands without performing discovery.
type SubcommandsReporter interface {
	// HasSubcommands returns true if the Command has subcommands or —
	// when discovery has not been performed — may have subcommands.
	HasSubcommands() bool
}

// HasSubcommands returns true if the given Command has subcommands or —
// when the Command defers discovery — may have subcommands.
//
// Unlike calling Subcommands() and checking its length, HasSubcommands
// never triggers discovery. Use it when a cheap, possibly-approximate
// answer is preferable to an exact, possibly-expensive one.
func HasSubcommands(command Command) bool {
	if r, ok := command.(SubcommandsReporter); ok {
		return r.HasSubcommands()
	}
	cmds, err := command.Subcommands()
	return err == nil && len(cmds) > 0
}
