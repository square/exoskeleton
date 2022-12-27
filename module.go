package exoskeleton

import (
	"os"
	"path/filepath"

	"github.com/squareup/exoskeleton/pkg/shellcomp"
)

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
	Subcommands() Commands
}

type module struct {
	executable
	cmds       Commands
	discoverer *discoverer
}

func (m *module) Exec(e *Entrypoint, args, env []string) error {
	return e.printModuleHelp(m, args)
}

func (m *module) Complete(_ *Entrypoint, args, _ []string) ([]string, shellcomp.Directive, error) {
	return m.Subcommands().completionsFor(args)
}

func (m *module) Summary() (string, error) {
	return getMessageFromDir(m.path, "summary")
}

func (m *module) Help() (string, error) {
	return getMessageFromDir(m.path, "help")
}

func (m *module) Subcommands() Commands {
	if m.cmds == nil {
		m.discoverer.discoverIn(filepath.Dir(m.path), m, &m.cmds)
	}

	return m.cmds
}

func getMessageFromDir(path string, flag string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return getMessageFromMagicComments(f, flag)
}
