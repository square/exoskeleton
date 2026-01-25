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
