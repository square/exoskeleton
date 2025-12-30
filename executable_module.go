package exoskeleton

import (
	"github.com/square/exoskeleton/pkg/shellcomp"
)

type executableModule struct {
	executableCommand
	cmds       Commands
	discoverer DiscoveryContext
}

func (m *executableModule) Summary() (string, error) {
	if m.cmds == nil {
		if err := m.discover(); err != nil {
			return "", err
		}
	}

	return *m.summary, nil
}

func (m *executableModule) Exec(e *Entrypoint, args, env []string) error {
	return e.printModuleHelp(m, args)
}

func (m *executableModule) Complete(_ *Entrypoint, args, _ []string) ([]string, shellcomp.Directive, error) {
	return completionsForModule(m, args)
}

func (m *executableModule) Subcommands() (Commands, error) {
	if m.cmds == nil {
		if err := m.discover(); err != nil {
			return Commands{}, err
		}
	}

	return m.cmds, nil
}

// discover invokes an executable with `--describe-commands` and constructs a tree
// of modules and subcommands (all to be invoked through the given executable)
// from the JSON output.
func (m *executableModule) discover() error {
	descriptor, err := describeCommands(m)
	if err != nil {
		return err
	}

	m.summary = descriptor.Summary
	m.cmds = toCommands(m, descriptor.Commands, nil, m.discoverer)
	return nil
}

func toCommands(parent *executableModule, descriptors []*commandDescriptor, args []string, d DiscoveryContext) Commands {
	cmds := Commands{}
	for _, descriptor := range descriptors {
		c := &executableCommand{
			parent:       parent,
			discoveredIn: parent.discoveredIn,
			path:         parent.path,
			args:         append(args, descriptor.Name),
			name:         descriptor.Name,
			summary:      descriptor.Summary,
			executor:     parent.executor,
		}

		if len(descriptor.Commands) > 0 && d.MaxDepth() != 0 {
			m := &executableModule{executableCommand: *c}
			m.cmds = toCommands(m, descriptor.Commands, append(args, m.name), d.Next())
			cmds = append(cmds, m)
		} else {
			cmds = append(cmds, c)
		}
	}
	return cmds
}
