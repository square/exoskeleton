package exoskeleton

import (
	"encoding/json"

	"github.com/square/exoskeleton/pkg/shellcomp"
)

type executableModule struct {
	executableCommand
	cmds       Commands
	discoverer discoverer
}

func (m *executableModule) Summary() (string, error) {
	if m.cmds == nil {
		m.discover()
	}

	return m.summary, nil
}

func (m *executableModule) Exec(e *Entrypoint, args, env []string) error {
	return e.printModuleHelp(m, args)
}

func (m *executableModule) Complete(_ *Entrypoint, args, _ []string) ([]string, shellcomp.Directive, error) {
	return m.Subcommands().completionsFor(args)
}

func (m *executableModule) Subcommands() Commands {
	if m.cmds == nil {
		m.discover()
	}

	return m.cmds
}

// discover invokes an executable with `--describe-commands` and constructs a tree
// of modules and subcommands (all to be invoked through the given executable)
// from the JSON output.
func (m *executableModule) discover() {
	cmd := m.Command("--describe-commands")
	cmd.Stderr = nil
	output, err := cmd.Output()
	if err != nil {
		m.discoverer.onError(DiscoveryError{Cause: err, Path: m.path})
		return
	}

	var descriptor *commandDescriptor
	if err := json.Unmarshal(output, &descriptor); err != nil {
		m.discoverer.onError(DiscoveryError{Cause: err, Path: m.path})
		return
	}

	m.name = descriptor.Name
	m.summary = descriptor.Summary
	m.cmds = m.discoverer.toCommands(m, descriptor.Commands, nil)
}

func (d *discoverer) toCommands(parent *executableModule, descriptors []*commandDescriptor, args []string) Commands {
	cmds := Commands{}
	for _, descriptor := range descriptors {
		c := &executableCommand{
			parent:       parent,
			discoveredIn: parent.discoveredIn,
			path:         parent.path,
			args:         append(args, descriptor.Name),
			name:         descriptor.Name,
			summary:      descriptor.Summary,
		}

		depth := d.depth + len(args)
		if len(descriptor.Commands) > 0 && (d.maxDepth == -1 || depth < d.maxDepth) {
			m := &executableModule{executableCommand: *c}
			m.cmds = d.toCommands(m, descriptor.Commands, append(args, m.name))
			cmds = append(cmds, m)
		} else {
			cmds = append(cmds, c)
		}
	}
	return cmds
}
