package exoskeleton

import (
	"io/fs"
	"os"
	"path/filepath"
)

type discoverer struct {
	maxDepth   int
	depth      int
	onError    func(error)
	modulefile string
}

func (e *Entrypoint) discoverIn(paths []string) (all Commands) {
	d := &discoverer{
		onError:    e.onError,
		modulefile: e.moduleMetadataFilename,
		maxDepth:   e.maxDepth,
	}
	for _, path := range paths {
		d.discoverIn(path, e, &all)
	}
	return
}

func (d *discoverer) discoverIn(path string, parent Module, all *Commands) {
	files, err := os.ReadDir(path)
	if err != nil {
		d.onError(DiscoveryError{Cause: err, Path: path})
		// No return. We may have a partial list of files: "ReadDir returns the entries
		// it was able to read before the error, along with the error"
	}

	for _, file := range files {
		name := file.Name()

		if file.Type()&fs.ModeSymlink != 0 {
			p := filepath.Join(path, name)
			file, err = followSymlinks(p)
			if err != nil {
				d.onError(DiscoveryError{Cause: err, Path: p})
				continue // skip this file
			}
		}

		if file.IsDir() {
			modulefilePath := filepath.Join(path, name, d.modulefile)

			// Don't search directories that exceed the configured maxDepth
			// or that don't contain the configured modulefile.
			if (d.maxDepth == -1 || d.depth < d.maxDepth) && exists(modulefilePath) {
				*all = append(*all, &directoryModule{
					executableCommand: executableCommand{
						parent:       parent,
						path:         modulefilePath,
						name:         name,
						discoveredIn: path,
					},
					discoverer: discoverer{
						maxDepth:   d.maxDepth,
						depth:      d.depth + 1,
						onError:    d.onError,
						modulefile: d.modulefile,
					},
				})
			}

		} else if ok, err := isExecutable(file); err != nil {
			d.onError(DiscoveryError{Cause: err, Path: path})
		} else if ok {
			*all = append(*all, &executableCommand{
				parent:       parent,
				path:         filepath.Join(path, name),
				name:         name,
				discoveredIn: path,
			})
		}
	}
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
