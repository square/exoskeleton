package exoskeleton

import (
	"path/filepath"

	"github.com/square/exoskeleton/pkg/shellcomp"
)

type directoryModule struct {
	executableCommand
	cmds       Commands
	discoverer discoverer
}

func (m *directoryModule) Exec(e *Entrypoint, args, env []string) error {
	return e.printModuleHelp(m, args)
}

func (m *directoryModule) Complete(_ *Entrypoint, args, _ []string) ([]string, shellcomp.Directive, error) {
	return m.Subcommands().completionsFor(args)
}

func (m *directoryModule) Summary() string {
	return m.entrypoint.readSummaryFromModulefile(m)
}

func (m *directoryModule) Help() string {
	panic("Unused")
}

func (m *directoryModule) Subcommands() Commands {
	if m.cmds == nil {
		m.discoverer.discoverIn(filepath.Dir(m.path), m, &m.cmds)
	}

	return m.cmds
}
