package exoskeleton

import (
	"io/fs"
	"path/filepath"
)

// StandaloneExecutableContract handles executable files that respond to --summary and --help.
//
// Note: This contract should be ordered AFTER ScriptCommandContract and ExecutableModuleContract
// in the contract list, as it matches all executable files.
type StandaloneExecutableContract struct{}

func (c *StandaloneExecutableContract) BuildCommand(path string, info fs.DirEntry, parent Command, d DiscoveryContext) (Command, error) {
	// Only applies to files
	if info.IsDir() {
		return nil, ErrNotApplicable
	}

	// Must be executable
	if ok, err := isExecutable(info); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrNotApplicable
	}

	return &executableCommand{
		parent:       parent,
		path:         path,
		name:         filepath.Base(path),
		discoveredIn: filepath.Dir(path),
		executor:     d.Executor(),
		cache:        d.Cache(),
	}, nil
}
