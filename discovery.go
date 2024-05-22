package exoskeleton

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const executableModuleExtension = ".exoskeleton"

type discoverer struct {
	maxDepth   int
	depth      int
	onError    func(error)
	modulefile string
}

func (e *Entrypoint) discoverIn(paths []string) Commands {
	all := Commands{}
	d := &discoverer{
		onError:    e.onError,
		modulefile: e.moduleMetadataFilename,
		maxDepth:   e.maxDepth,
	}
	for _, path := range paths {
		d.discoverIn(path, e, &all)
	}
	return all
}

func (d *discoverer) discoverIn(path string, parent Module, all *Commands) {
	files, err := os.ReadDir(path)
	if err != nil {
		d.onError(DiscoveryError{Cause: err, Path: path})
		// No return. We may have a partial list of files: "ReadDir returns the entries
		// it was able to read before the error, along with the error"
	}

	for _, file := range files {
		if cmd, err := d.buildCommand(path, parent, file); err != nil {
			d.onError(err)
		} else if cmd != nil {
			*all = append(*all, cmd)
		}
	}
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

	if file.IsDir() {
		modulefilePath := filepath.Join(path, d.modulefile)

		// If the directory doesn't contain the modulefile, it is just a regular directory
		if !exists(modulefilePath) {
			return nil, nil
		}

		// Stop discovering modules if we've searched past maxDepth
		if d.maxDepth >= 0 && d.depth >= d.maxDepth {
			return nil, nil
		}

		return &directoryModule{
			executableCommand: executableCommand{
				parent:       parent,
				path:         modulefilePath,
				name:         name,
				discoveredIn: discoveredIn,
			},
			discoverer: discoverer{
				maxDepth:   d.maxDepth,
				depth:      d.depth + 1,
				onError:    d.onError,
				modulefile: d.modulefile,
			},
		}, nil
	} else {
		if ok, err := isExecutable(file); err != nil {
			return nil, DiscoveryError{Cause: err, Path: path}
		} else if !ok {
			// If the file isn't executable, it is just a regular file
			return nil, nil
		}

		executable := &executableCommand{
			parent:       parent,
			path:         path,
			name:         strings.TrimSuffix(name, executableModuleExtension),
			discoveredIn: discoveredIn,
		}

		// if the executable has the extension ".exoskeleton" then we should treat it as a module.
		if filepath.Ext(name) == executableModuleExtension && (d.maxDepth == -1 || d.depth < d.maxDepth) {
			// Execute this command with the `--describe-commands` flag to get the subcommands
			if module, err := d.discoverSubcommands(executable); err != nil {
				return nil, DiscoveryError{Cause: err, Path: path}
			} else {
				return module, nil
			}
		}

		return executable, nil
	}
}

type commandDescriptor struct {
	Name     string               `json:"name"`
	Summary  string               `json:"summary"`
	Commands []*commandDescriptor `json:"commands,omitempty"`
}

// discoverSubcommands invokes an executable with `--describe-commands` and constructs
// a tree of modules and subcommands (all to be invoked through the given executable)
// from the JSON output.
func (d *discoverer) discoverSubcommands(executable *executableCommand) (Command, error) {
	cmd := executable.Command("--describe-commands")
	cmd.Stderr = nil
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var descriptor *commandDescriptor
	if err := json.Unmarshal(output, &descriptor); err != nil {
		return nil, err
	}

	executable.name = descriptor.Name
	executable.summary = descriptor.Summary
	m := &executableModule{executableCommand: *executable}
	m.cmds = d.toCommands(m, descriptor.Commands, nil)
	return m, nil
}

func (d *discoverer) toCommands(parent *executableModule, descriptors []*commandDescriptor, args []string) Commands {
	var cmds Commands
	for _, descriptor := range descriptors {
		c := &executableCommand{
			parent:       parent,
			discoveredIn: parent.discoveredIn,
			path:         parent.path,
			args:         append(args, descriptor.Name),
			name:         descriptor.Name,
			summary:      descriptor.Summary,
		}

		depth := d.depth + len(args) + 1
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
