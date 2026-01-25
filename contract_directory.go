package exoskeleton

import (
	"io/fs"
	"os"
	"path/filepath"
)

// DirectoryContract handles directories containing a module metadata file.
//
// The module metadata filename is configured on the Entrypoint and passed through
// the discoverer, not stored on the contract itself.
type DirectoryContract struct {
	MetadataFilename string
}

func (c *DirectoryContract) BuildCommand(path string, info fs.DirEntry, parent Command, d DiscoveryContext) (Command, error) {
	// Only applies to directories
	if !info.IsDir() {
		return nil, ErrNotApplicable
	}

	modulefilePath := filepath.Join(path, c.MetadataFilename)

	// If the directory doesn't contain the modulefile, it's just a regular directory
	if !exists(modulefilePath) {
		return nil, ErrNotApplicable
	}

	// Stop discovering modules if we've searched past maxDepth
	if d.MaxDepth() == 0 {
		return nil, nil // Ignore due to depth limit
	}

	return &directoryCommand{
		executableCommand: executableCommand{
			parent:       parent,
			path:         modulefilePath,
			name:         filepath.Base(path),
			discoveredIn: filepath.Dir(path),
			executor:     d.Executor(),
			cache:        d.Cache(),
		},
		discoverer: d.Next(),
	}, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
