package exoskeleton

import (
	"github.com/square/exoskeleton/pkg/shellcomp"
)

type builtinModule struct {
	parent      Module
	definition  *EmbeddedModule
	subcommands Commands
}

func (m *builtinModule) Parent() Module           { return m.parent }
func (m *builtinModule) Path() string             { return m.parent.Path() }
func (m *builtinModule) Name() string             { return m.definition.Name }
func (m *builtinModule) Summary() (string, error) { return m.definition.Summary, nil }
func (m *builtinModule) Help() (string, error)    { return "", nil }
func (m *builtinModule) Subcommands() Commands    { return m.subcommands }

func (m *builtinModule) Exec(e *Entrypoint, args, env []string) error {
	return e.printModuleHelp(m, args)
}

func (m *builtinModule) Complete(e *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error) {
	return m.Subcommands().completionsFor(args)
}
