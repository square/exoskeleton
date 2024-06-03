package exoskeleton

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/square/exoskeleton/pkg/shellcomp"
)

// CommandNotFoundFunc is a function that accepts a Null Command object. It is called when a command is not found.
type CommandNotFoundFunc func(*Entrypoint, Command)

// ErrorFunc is called when an error occurs.
type ErrorFunc func(*Entrypoint, error)

// MenuHeadingForFunc is a function that is expected to return the heading
// under which a command should be listed when it is printed in a menu.
// (The default value is "COMMANDS".)
//
// It accepts the Module whose subcommands are being listed and the Command
// whose heading should be returned.
type MenuHeadingForFunc func(Module, Command) string

// Entrypoint is the root of an exoskeleton CLI application.
type Entrypoint struct {
	path                     string
	name                     string
	cmds                     Commands
	maxDepth                 int
	cachePath                string
	menuHeadingFor           MenuHeadingForFunc
	moduleMetadataFilename   string
	completeHelp             string
	helpHelp                 string
	whichHelp                string
	errorCallbacks           []ErrorFunc
	commandNotFoundCallbacks []CommandNotFoundFunc
	cmdsToAppend             []Command
	cmdsToPrepend            []Command
}

func (e *Entrypoint) Parent() Module        { return nil }
func (e *Entrypoint) Path() string          { return e.path }
func (e *Entrypoint) Name() string          { return e.name }
func (e *Entrypoint) Summary() string       { panic("Unused") }
func (e *Entrypoint) Help() string          { panic("Unused") }
func (e *Entrypoint) Subcommands() Commands { return e.cmds }

// New searches the given paths and constructs an Entrypoint with a list of commands
// discovered in those paths. It also accepts options that can be used to customize
// the behavior of the Entrypoint.
func New(paths []string, options ...Option) (*Entrypoint, error) {
	path, err := os.Executable()
	if err != nil {
		return nil, err
	}

	self := newWithDefaults(path)

	options =
		append(
			options,
			PrependCommands(
				&EmbeddedCommand{
					Name:     "help",
					Help:     fmt.Sprintf(self.helpHelp, self.Name()),
					Exec:     helpExec,
					Complete: CompleteCommands,
				},
				&EmbeddedCommand{
					Name:     "which",
					Help:     fmt.Sprintf(self.whichHelp, self.Name()),
					Exec:     whichExec,
					Complete: CompleteCommands,
				},
				&EmbeddedCommand{
					Name:     "complete",
					Help:     fmt.Sprintf(self.completeHelp, self.Name()),
					Exec:     completeExec,
					Complete: nil,
				},
			),
		)

	for _, op := range options {
		op.Apply(self)
	}

	self.cmds =
		append(
			self.cmdsToPrepend,
			append(
				self.discoverIn(paths),
				self.cmdsToAppend...,
			)...,
		)

	return self, nil
}

func newWithDefaults(path string) *Entrypoint {
	name := filepath.Base(path)

	var cachePath string
	if dir, err := os.UserCacheDir(); err == nil {
		cachePath = filepath.Join(dir, name+".json")
	}

	return &Entrypoint{
		path:                   path,
		name:                   name,
		maxDepth:               -1,
		cachePath:              cachePath,
		menuHeadingFor:         func(_ Module, _ Command) string { return "COMMANDS" },
		moduleMetadataFilename: ".exoskeleton",
		completeHelp:           CompleteHelp,
		helpHelp:               HelpHelp,
		whichHelp:              WhichHelp,
		cmdsToPrepend:          []Command{},
		cmdsToAppend:           []Command{},
	}
}

func (e *Entrypoint) onError(err error) {
	for _, callback := range e.errorCallbacks {
		callback(e, err)
	}
}

func (e *Entrypoint) commandNotFound(cmd nullCommand) {
	for _, callback := range e.commandNotFoundCallbacks {
		callback(e, cmd)
	}

	usage := UsageRelativeTo(cmd, e)
	fmt.Fprintf(os.Stderr, "%s: no such command %s\n", e.Name(), usage)

	if suggestions := e.suggestionsFor(usage); len(suggestions) > 0 {
		fmt.Fprintln(os.Stderr, "Did you mean?")
		for _, suggestion := range suggestions {
			fmt.Fprintf(os.Stderr, "   %s\n", Usage(suggestion))
		}
	}
}

func (e *Entrypoint) Exec(_ *Entrypoint, rawArgs, env []string) error {
	return e.printModuleHelp(e, rawArgs)
}

func (e *Entrypoint) Complete(_ *Entrypoint, args, _ []string) (completions []string, directive shellcomp.Directive, err error) {
	completions, directive, err = e.Subcommands().completionsFor(args)
	completions = without(completions, `complete`) // Don't suggest the `complete` command
	return
}
