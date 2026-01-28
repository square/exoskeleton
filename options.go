package exoskeleton

import (
	"text/template"

	"github.com/square/exoskeleton/pkg/shellcomp"
)

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
	Commands []*EmbeddedCommand
}

// Apply invokes the optionFunc with the given Entrypoint.
func (fn optionFunc) Apply(e *Entrypoint) {
	fn(e)
}

// AppendCommands adds new embedded commands to the Entrypoint. The commands are added to
// the end of the list and will have the lowest precedence: an executable with the same name
// as one of these commands would override it.
func AppendCommands(cmds ...*EmbeddedCommand) Option {
	return (optionFunc)(func(e *Entrypoint) {
		e.cmdsToAppend = append(e.cmdsToAppend, buildCommands(e, cmds)...)
	})
}

// PrependCommands adds new embedded commands to the Entrypoint. The command are added to
// the front of the list and will have the highest precedence: an executable with the same name
// as one of these commands would be overridden by it.
func PrependCommands(cmds ...*EmbeddedCommand) Option {
	return (optionFunc)(func(e *Entrypoint) {
		e.cmdsToPrepend = append(buildCommands(e, cmds), e.cmdsToPrepend...)
	})
}

func buildCommands(parent Command, cmds []*EmbeddedCommand) []Command {
	var result []Command

	for _, cmd := range cmds {
		bc := &builtinCommand{parent: parent, definition: cmd}
		if len(cmd.Commands) > 0 {
			bc.subcommands = buildCommands(bc, cmd.Commands)
		}
		result = append(result, bc)
	}

	return result
}

// OnError registers a callback (ErrorFunc) to be invoked when a nonfatal error occurs.
//
// These are recoverable errors such as
//   - a broken symlink is encountered in one of the paths being searched
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

// AfterIdentify registers a callback (AfterIdentifyFunc) to be invoked with any command
// that is successfully identified.
func AfterIdentify(fn AfterIdentifyFunc) Option {
	return (optionFunc)(func(e *Entrypoint) { e.afterIdentifyCallbacks = append(e.afterIdentifyCallbacks, fn) })
}

// WithName sets the name of the entrypoint.
// (By default, this is the basename of the executable.)
func WithName(value string) Option {
	return (optionFunc)(func(e *Entrypoint) { e.name = value })
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

// WithMenuTemplate sets the template that will be used to render help for modules.
// The template will be executed with an instance of exoskeleton.Menu as its data.
func WithMenuTemplate(value *template.Template) Option {
	return (optionFunc)(func(e *Entrypoint) { e.menuTemplate = value })
}

// WithExecutor supplies a function that executes a subcommand.
// The default executor calls `Run()` on the command and returns the error.
func WithExecutor(value ExecutorFunc) Option {
	return (optionFunc)(func(e *Entrypoint) { e.executor = value })
}

// WithModuleMetadataFilename sets the filename to use for module metadata.
// (Default: ".exoskeleton")
func WithModuleMetadataFilename(value string) Option {
	return (optionFunc)(func(e *Entrypoint) { e.moduleMetadataFilename = value })
}

// WithContracts sets the contracts used during discovery.
// Contracts are tried in order; the first that doesn't return ErrNotApplicable builds the command.
//
// The default contracts are:
//  1. DirectoryContract (directories that contain the module metadata file)
//  2. ExecutableContract (executables with .exoskeleton extension which must implement --describe-commands)
//  3. ShellScriptContract (shell scripts with magic comments)
//  4. StandaloneExecutableContract (all other executables which must implement --summary)
func WithContracts(contracts ...Contract) Option {
	return (optionFunc)(func(e *Entrypoint) { e.contracts = contracts })
}

// WithCache sets a cache for expensive operations like command execution.
// If not set, no caching is performed.
func WithCache(c Cache) Option {
	return (optionFunc)(func(e *Entrypoint) { e.cache = c })
}
