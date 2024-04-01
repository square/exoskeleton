package exoskeleton

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type discoverer struct {
	entrypoint *Entrypoint
	depth      int
}

func (e *Entrypoint) discoverIn(paths []string) (all Commands) {
	d := &discoverer{entrypoint: e}
	for _, path := range paths {
		d.discoverIn(path, e, &all)
	}
	return
}

func (d *discoverer) discoverIn(p string, parent Module, all *Commands) {
	files, err := os.ReadDir(p)
	if err != nil {
		d.onError(DiscoveryError{Cause: err, Path: p})
		// No return. We may have a partial list of files: "ReadDir returns the entries
		// it was able to read before the error, along with the error"
	}

	for _, file := range files {
		name := file.Name()

		if file.Type()&fs.ModeSymlink != 0 {
			p := path.Join(p, name)
			file, err = followSymlinks(p)
			if err != nil {
				d.onError(DiscoveryError{Cause: err, Path: p})
				continue // skip this file
			}
		}

		if file.IsDir() {
			modulefilePath := path.Join(p, name, d.entrypoint.moduleMetadataFilename)

			// Don't search directories that exceed the configured maxDepth
			// or that don't contain the configured modulefile.
			if (d.entrypoint.maxDepth == -1 || d.depth < d.entrypoint.maxDepth) && exists(modulefilePath) {
				*all = append(*all, &directoryModule{
					executableCommand: executableCommand{
						entrypoint:   d.entrypoint,
						parent:       parent,
						path:         modulefilePath,
						name:         name,
						discoveredIn: p,
					},
					discoverer: &discoverer{
						entrypoint: d.entrypoint,
						depth:      d.depth + 1,
					},
				})
			}

		} else if ok, err := isExecutable(file); err != nil {
			d.onError(DiscoveryError{Cause: err, Path: p})
		} else if ok {
			*all = append(*all, &executableCommand{
				entrypoint:   d.entrypoint,
				parent:       parent,
				path:         path.Join(p, name),
				name:         name,
				discoveredIn: p,
			})
		}
	}
}

func (d *discoverer) onError(err error) {
	d.entrypoint.onError(err)
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

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
