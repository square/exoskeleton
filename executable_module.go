package exoskeleton

import (
	"github.com/square/exoskeleton/pkg/shellcomp"
)

type executableModule struct {
	executableCommand
	cmds Commands
}

func (m *executableModule) Exec(e *Entrypoint, args, env []string) error {
	return e.printModuleHelp(m, args)
}

func (m *executableModule) Complete(_ *Entrypoint, args, _ []string) ([]string, shellcomp.Directive, error) {
	return m.Subcommands().completionsFor(args)
}

func (m *executableModule) Subcommands() Commands { return m.cmds }
