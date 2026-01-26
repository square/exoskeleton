package exoskeleton

import (
	"io/fs"
	"os"
	"path/filepath"
)

// ShellScriptContract handles shell scripts with magic comments.
//
// Scripts are detected by checking for a shebang (#!) at the start of the file.
// They provide metadata via magic comments like "# SUMMARY:" and "# HELP:".
type ShellScriptContract struct{}

func (c *ShellScriptContract) BuildCommand(path string, info fs.DirEntry, parent Module, d DiscoveryContext) (Command, error) {
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

	// Must start with shebang
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buffer := make([]byte, 2)
	if _, err := f.Read(buffer); err != nil {
		return nil, err
	}
	if string(buffer) != "#!" {
		return nil, ErrNotApplicable // Not a shell script
	}

	return &shellScriptCommand{
		executableCommand: executableCommand{
			parent:       parent,
			path:         path,
			name:         filepath.Base(path),
			discoveredIn: filepath.Dir(path),
			executor:     d.Executor(),
			cache:        d.Cache(),
		},
	}, nil
}

// shellScriptCommand implements the Command interface for shell scripts that use magic comments.
// It extends executableCommand but overrides Summary() and Help() to read magic comments
// instead of executing with flags.
type shellScriptCommand struct {
	executableCommand
}

func (cmd *shellScriptCommand) Summary() (string, error) {
	if cmd.summary != nil {
		return *cmd.summary, nil
	}

	return cmd.cache.Fetch(cmd, "summary", func() (string, error) {
		return readSummaryFromShellScript(cmd)
	})
}

func (cmd *shellScriptCommand) Help() (string, error) {
	return readHelpFromShellScript(cmd)
}
