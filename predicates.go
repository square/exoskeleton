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

// IsHidden returns true if the given Command asks not to be listed in menus,
// completions, or suggestions. Hidden commands can still be identified and run.
func IsHidden(command Command) bool {
	d, ok := command.(OpenCLIDescriber)
	if !ok {
		return false
	}
	c, err := d.OpenCLICommand()
	return err == nil && c != nil && c.Hidden
}
