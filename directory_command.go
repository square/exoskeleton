package exoskeleton

import (
	"path/filepath"

	"github.com/square/exoskeleton/v2/pkg/shellcomp"
)

type directoryCommand struct {
	executableCommand
	cmds       Commands
	discoverer DiscoveryContext
}

func (m *directoryCommand) Exec(e *Entrypoint, args, env []string) error {
	return e.printModuleHelp(m, args)
}

func (m *directoryCommand) Complete(_ *Entrypoint, args, _ []string) ([]string, shellcomp.Directive, error) {
	return completionsForSubcommands(m, args)
}

func (m *directoryCommand) Summary() (string, error) {
	return m.cache.Fetch(m, "summary", func() (string, error) {
		return readSummaryFromModulefile(m)
	})
}

func (m *directoryCommand) Help() (string, error) {
	panic("Unused")
}

func (m *directoryCommand) Subcommands() (Commands, error) {
	if m.cmds == nil && m.discoverer != nil {
		m.cmds, _ = m.discoverer.DiscoverIn(filepath.Dir(m.path), m)
		// TODO: return errors from discovery
	}

	return m.cmds, nil
}
