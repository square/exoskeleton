package exoskeleton

import (
	"io/fs"
	"path/filepath"
	"strings"
)

const executableModuleExtension = ".exoskeleton"

// ExecutableContract handles executables that respond to --describe-commands.
//
// The extension is hardcoded to ".exoskeleton", and the --describe-commands flag is also hardcoded.
type ExecutableContract struct {
}

func (c *ExecutableContract) BuildCommand(path string, info fs.DirEntry, parent Module, d DiscoveryContext) (Command, error) {
	// Only applies to files
	if info.IsDir() {
		return nil, ErrNotApplicable
	}

	// Must have the configured extension
	name := filepath.Base(path)
	if filepath.Ext(name) != executableModuleExtension {
		return nil, ErrNotApplicable
	}

	// Must be executable
	if ok, err := isExecutable(info); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrNotApplicable
	}

	commandName := strings.TrimSuffix(name, executableModuleExtension)

	// Stop discovering modules if we've searched past maxDepth
	// But still create a regular executableCommand (not a module)
	if d.MaxDepth() == 0 {
		return &executableCommand{
			parent:       parent,
			path:         path,
			name:         commandName,
			discoveredIn: filepath.Dir(path),
			executor:     d.Executor(),
			cache:        d.Cache(),
		}, nil
	}

	return &executableModule{
		executableCommand: executableCommand{
			parent:       parent,
			path:         path,
			name:         commandName,
			discoveredIn: filepath.Dir(path),
			executor:     d.Executor(),
			cache:        d.Cache(),
		},
		discoverer: d.Next(),
	}, nil
}
