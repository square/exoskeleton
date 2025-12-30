package exoskeleton

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type discoverer struct {
	maxDepth  int
	onError   func(error)
	executor  ExecutorFunc
	contracts []Contract
}

type DiscoveryContext interface {
	DiscoverIn(path string, parent Module) (Commands, []error)
	Executor() ExecutorFunc
	MaxDepth() int
	Next() DiscoveryContext
}

func (d *discoverer) Next() DiscoveryContext {
	return &discoverer{
		maxDepth:  d.MaxDepth() - 1,
		onError:   d.onError,
		executor:  d.executor,
		contracts: d.contracts,
	}
}

var _ DiscoveryContext = &discoverer{}

func (e *Entrypoint) discoverIn(paths []string) Commands {
	all := Commands{}
	d := &discoverer{
		onError:   e.onError,
		executor:  e.executor,
		maxDepth:  e.maxDepth,
		contracts: e.contracts,
	}
	for _, path := range paths {
		cmds, _ := d.DiscoverIn(path, e)
		all = append(all, cmds...)
	}
	return all
}

func (d *discoverer) Executor() ExecutorFunc { return d.executor }
func (d *discoverer) MaxDepth() int          { return d.maxDepth }

func (d *discoverer) DiscoverIn(path string, parent Module) (Commands, []error) {
	var all Commands
	var errs []error

	files, err := os.ReadDir(path)
	if err != nil {
		if d.onError != nil {
			d.onError(err)
		}
		errs = append(errs, DiscoveryError{Cause: err, Path: path})
		// No return. We may have a partial list of files: "ReadDir returns the entries
		// it was able to read before the error, along with the error"
	}

	for _, file := range files {
		if cmd, err := d.buildCommand(path, parent, file); err != nil {
			if d.onError != nil {
				d.onError(err)
			}
			errs = append(errs, err)
		} else if cmd != nil {
			all = append(all, cmd)
		}
	}

	return all, errs
}

func (d *discoverer) buildCommand(discoveredIn string, parent Module, file fs.DirEntry) (Command, error) {
	name := file.Name()
	path := filepath.Join(discoveredIn, name)

	var err error
	if file.Type()&fs.ModeSymlink != 0 {
		file, err = followSymlinks(path)
		if err != nil {
			return nil, DiscoveryError{Cause: err, Path: path}
		}
	}

	// Try each contract in order
	for _, contract := range d.contracts {
		if cmd, err := contract.BuildCommand(path, file, parent, d); err == nil {
			return cmd, nil
		} else if !errors.Is(err, ErrNotApplicable) {
			return nil, DiscoveryError{Cause: err, Path: path} // Contract failed with real error
		}
		// Contract doesn't apply, try next one
	}

	// No contract applies. File is ignored.
	return nil, nil
}

func followSymlinks(path string) (fs.DirEntry, error) {
	if realPath, err := filepath.EvalSymlinks(path); err != nil {
		return nil, SymlinkError{Cause: err, Path: path}
	} else if info, err := os.Lstat(realPath); err != nil {
		return nil, SymlinkError{Cause: err, Path: path}
	} else {
		return fs.FileInfoToDirEntry(info), nil
	}
}

func isExecutable(file fs.DirEntry) (bool, error) {
	if info, err := file.Info(); err != nil {
		return false, err
	} else {
		return info.Mode()&0111 != 0, nil // is x bit set for User, Group, or Other
	}
}
