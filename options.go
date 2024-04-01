package exoskeleton

import "github.com/square/exoskeleton/pkg/shellcomp"

// An Option applies optional changes to an Exoskeleton Entrypoint.
type Option interface {
	Apply(*Entrypoint)
}

// optionFunc is a function that adheres to the Option interface.
type optionFunc func(*Entrypoint)

// ExecFunc is called when an built-in command is run.
type ExecFunc func(e *Entrypoint, args, env []string) error

// CompleteFunc is called when an built-in command is asked to supply shell completions.
type CompleteFunc func(e *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error)

// EmbeddedCommand defines a built-in command that can be added to an Entrypoint
// (as opposed to an executable external to the Entrypoint).
type EmbeddedCommand struct {
	Name     string
	Summary  string
	Help     string
	Exec     ExecFunc
	Complete CompleteFunc
}

// EmbeddedCommand defines a built-in module that can be added to an Entrypoint
// (as opposed to a directory external to the Entrypoint).
type EmbeddedModule struct {
	Name     string
	Summary  string
	Commands []interface{}
}

// Apply invokes the optionFunc with the given Entrypoint.
func (fn optionFunc) Apply(e *Entrypoint) {
	fn(e)
}

// AppendCommands adds new embedded commands to the Entrypoint. The commands are added to
// the end of the list and will have the lowest precedence: an executable with the same name
// as one of these commands would override it.
func AppendCommands(cmds ...interface{}) Option {
	return (optionFunc)(func(e *Entrypoint) {
		e.cmdsToAppend = append(e.cmdsToAppend, buildCommands(e, cmds)...)
	})
}

// PrependCommands adds new embedded commands to the Entrypoint. The command are added to
// the front of the list and will have the highest precedence: an executable with the same name
// as one of these commands would be overridden by it.
func PrependCommands(cmds ...interface{}) Option {
	return (optionFunc)(func(e *Entrypoint) {
		e.cmdsToPrepend = append(buildCommands(e, cmds), e.cmdsToPrepend...)
	})
}

func buildCommands(m Module, cmds []interface{}) []Command {
	var result []Command

	for _, cmd := range cmds {
		switch v := cmd.(type) {
		case *EmbeddedCommand:
			result = append(result, &builtinCommand{m, v})
		case *EmbeddedModule:
			module := &builtinModule{m, v, nil}
			module.subcommands = buildCommands(module, v.Commands)
			result = append(result, module)
		default:
			panic("Invalid command type")
		}
	}

	return result
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
