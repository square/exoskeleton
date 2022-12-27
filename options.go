package exoskeleton

// An Option applies optional changes to an Exoskeleton Entrypoint.
type Option interface {
	Apply(*Entrypoint)
}

// optionFunc is a function that adheres to the Option interface.
type optionFunc func(*Entrypoint)

// Apply invokes the optionFunc with the given Entrypoint.
func (fn optionFunc) Apply(e *Entrypoint) {
	fn(e)
}

// AppendCommand adds a new built-in command to the Entrypoint. The command is added to
// the end of the list and will have the lowest precedence: an executable with the same name
// would override it.
func AppendCommand(name, summary, help string, exec ExecFunc, complete CompleteFunc) Option {
	return (optionFunc)(func(e *Entrypoint) {
		e.cmdsToAppend = append(e.cmdsToAppend, &builtinCommand{e, name, summary, help, exec, complete})
	})
}

// PrependCommand adds a new built-in command to the Entrypoint. The command is added to
// the front of the list and will have the highest precedence: an executable with the same name
// would be overridden by it.
func PrependCommand(name, summary, help string, exec ExecFunc, complete CompleteFunc) Option {
	return (optionFunc)(func(e *Entrypoint) {
		e.cmdsToPrepend = append([]Command{&builtinCommand{e, name, summary, help, exec, complete}}, e.cmdsToPrepend...)
	})
}

// OnError registers a callback (ErrorFunc) to be invoked when a nonfatal error occurs.
//
// These are recoverable errors such as
//   - a broken symlink is encountered in one of the paths being searched
//   - exoskeleton is unable to leverage its cache because it is unable to read or write to it
//   - a command exits unnsuccessfully when invoked with --summary or --help
func OnError(fn ErrorFunc) Option {
	return (optionFunc)(func(e *Entrypoint) { e.errorCallbacks = append(e.errorCallbacks, fn) })
}

// OnCommandNotFound registers a callback (CommandNotFoundFunc) to be invoked when
// the command a user attempted to execute is not found. The callback is also invoked when
// the user asks for help on a command that can not be found.
func OnCommandNotFound(fn CommandNotFoundFunc) Option {
	return (optionFunc)(func(e *Entrypoint) { e.commandNotFoundCallbacks = append(e.commandNotFoundCallbacks, fn) })
}

// WithMaxDepth sets the maximum depth of the command tree.
//
// A value of 0 prohibits any submodules. All subcommands are leaves of the Entrypoint.
// (i.e. If the Go CLI were an exoskeleton, 'go doc' would be allowed, 'go mod tidy' would not.)
//
// A value of -1 (the default value) means there is no maximum depth.
func WithMaxDepth(value int) Option {
	return (optionFunc)(func(e *Entrypoint) { e.maxDepth = value })
}

// WithMenuHeadingFor allows you to supply a function that determines the heading
// a Command should be listed under in the main menu.
func WithMenuHeadingFor(fn MenuHeadingForFunc) Option {
	return (optionFunc)(func(e *Entrypoint) { e.menuHeadingFor = fn })
}

// WithModuleMetadataFilename sets the filename to use for module metadata.
// (Default: ".exoskeleton")
func WithModuleMetadataFilename(value string) Option {
	return (optionFunc)(func(e *Entrypoint) { e.moduleMetadataFilename = value })
}

// WithHelpText allows you to override the help text for a built-in command.
// There are three build-in commands: 'complete', 'help', and 'which'.
func WithHelpText(command string, value string) Option {
	switch command {
	case `complete`:
		return (optionFunc)(func(e *Entrypoint) { e.completeHelp = value })
	case `help`:
		return (optionFunc)(func(e *Entrypoint) { e.helpHelp = value })
	case `which`:
		return (optionFunc)(func(e *Entrypoint) { e.whichHelp = value })
	default:
		panic(`Invalid command: ` + command)
	}
}
