package exoskeleton

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/block/opencli-go"
	"github.com/square/exit"
)

// OpenCLIContract handles executables that respond to --help-opencli
// with an OpenCLI document (github.com/block/opencli-go).
//
// Unlike ExecutableContract, this contract does not require any particular
// file extension. Any executable file is eligible.
type OpenCLIContract struct{}

func (c *OpenCLIContract) BuildCommand(path string, info fs.DirEntry, parent Command, d DiscoveryContext) (Command, error) {
	if info.IsDir() {
		return nil, ErrNotApplicable
	}

	if ok, err := isExecutable(info); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrNotApplicable
	}

	name := filepath.Base(path)

	if d.MaxDepth() == 0 {
		return &executableCommand{
			parent:       parent,
			path:         path,
			name:         name,
			discoveredIn: filepath.Dir(path),
			executor:     d.Executor(),
			cache:        d.Cache(),
		}, nil
	}

	return &executableCommand{
		parent:       parent,
		path:         path,
		name:         name,
		discoveredIn: filepath.Dir(path),
		executor:     d.Executor(),
		cache:        d.Cache(),
		discoverer:   d.Next(),
		describe:     describeOpenCLI,
	}, nil
}

// describeOpenCLI is a describeFunc that invokes --help-opencli
// and parses the OpenCLI JSON output into a commandDescriptor.
func describeOpenCLI(cmd *executableCommand) (*commandDescriptor, error) {
	out, err := cmd.cache.Fetch(cmd, "help-opencli", func() (string, error) {
		return helpOpenCLIRaw(cmd)
	})
	if err != nil {
		return nil, err
	}
	return parseOpenCLI(cmd, out)
}

// helpOpenCLIRaw executes --help-opencli and returns the raw JSON output.
func helpOpenCLIRaw(m *executableCommand) (string, error) {
	cmd := m.Command("--help-opencli")
	out, err := m.output(cmd)
	if err != nil {
		err = fmt.Errorf("exec '%s': %w", strings.Join(cmd.Args, " "), err)
		return "", exit.Wrap(
			CommandDescribeError{
				CommandError{
					Message: err.Error(),
					Command: m,
					Cause:   err,
				},
			},
			exit.InternalError,
		)
	}
	return string(out), nil
}

// parseOpenCLI parses OpenCLI JSON output into a commandDescriptor.
func parseOpenCLI(cmd *executableCommand, out string) (*commandDescriptor, error) {
	var doc opencli.Document
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		return &commandDescriptor{},
			exit.Wrap(
				CommandDescribeError{
					CommandError{
						Message: fmt.Sprintf("error parsing output from `%s --help-opencli`: %s", cmd.path, err),
						Command: cmd,
						Cause:   err,
					},
				},
				exit.InternalError,
			)
	}
	return opencliToDescriptor(doc.Command), nil
}

// opencliToDescriptor converts an opencli.Command into a commandDescriptor,
// mapping the fields that Exoskeleton uses and ignoring the rest.
func opencliToDescriptor(cmd opencli.Command) *commandDescriptor {
	d := &commandDescriptor{
		Name:    cmd.Name,
		Summary: cmd.Summary,
		Aliases: cmd.Aliases,
	}
	if cmd.DefaultCommand != nil {
		d.DefaultCommand = *cmd.DefaultCommand
	}
	if len(cmd.Commands) > 0 {
		d.Commands = make([]*commandDescriptor, len(cmd.Commands))
		for i, sub := range cmd.Commands {
			d.Commands[i] = opencliToDescriptor(sub)
		}
	}
	return d
}
