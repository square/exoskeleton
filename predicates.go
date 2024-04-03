package exoskeleton

// IsEmbedded returns true if the given Command is built into the exoskeleton.
func IsEmbedded(command Command) bool {
	_, isCmd := command.(*builtinCommand)
	_, isMod := command.(*builtinModule)
	return isCmd || isMod
}

// IsNull returns true if the given Command is a NullCommand and false if it is not.
func IsNull(command Command) bool {
	_, ok := command.(nullCommand)
	return ok
}

// IsModule returns true if the given Command is a Module and false if it is not.
func IsModule(command Command) bool {
	_, ok := command.(Module)
	return ok
}
