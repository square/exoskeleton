package exoskeleton

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenCLIContractRejectsDirectories(t *testing.T) {
	contract := &OpenCLIContract{}
	path := filepath.Join(fixtures, "nested-1")
	info, err := os.Lstat(path)
	assert.NoError(t, err)

	_, err = contract.BuildCommand(path, fs.FileInfoToDirEntry(info), nil, &discoverer{maxDepth: -1})
	assert.ErrorIs(t, err, ErrNotApplicable)
}

func TestOpenCLIContractRejectsNonExecutableFiles(t *testing.T) {
	contract := &OpenCLIContract{}
	path := filepath.Join(fixtures, "nested-1", ".exoskeleton")
	info, err := os.Lstat(path)
	assert.NoError(t, err)

	_, err = contract.BuildCommand(path, fs.FileInfoToDirEntry(info), nil, &discoverer{maxDepth: -1})
	assert.ErrorIs(t, err, ErrNotApplicable)
}

func TestOpenCLIContractBuildsCommandForAnyExecutable(t *testing.T) {
	contract := &OpenCLIContract{}
	path := filepath.Join(fixtures, "opencli-tool")
	info, err := os.Lstat(path)
	assert.NoError(t, err)

	d := &discoverer{maxDepth: -1, executor: defaultExecutor}
	cmd, err := contract.BuildCommand(path, fs.FileInfoToDirEntry(info), nil, d)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "opencli-tool", cmd.Name())
}

func TestOpenCLIContractUsesFilenameAsCommandName(t *testing.T) {
	contract := &OpenCLIContract{}
	path := filepath.Join(fixtures, "opencli-tool")
	info, err := os.Lstat(path)
	assert.NoError(t, err)

	d := &discoverer{maxDepth: -1, executor: defaultExecutor}
	cmd, err := contract.BuildCommand(path, fs.FileInfoToDirEntry(info), nil, d)
	assert.NoError(t, err)
	assert.Equal(t, "opencli-tool", cmd.Name())
}

func TestOpenCLIContractMaxDepthZeroCreatesLeaf(t *testing.T) {
	contract := &OpenCLIContract{}
	path := filepath.Join(fixtures, "opencli-tool")
	info, err := os.Lstat(path)
	assert.NoError(t, err)

	d := &discoverer{maxDepth: 0, executor: defaultExecutor}
	cmd, err := contract.BuildCommand(path, fs.FileInfoToDirEntry(info), nil, d)
	assert.NoError(t, err)

	// At maxDepth=0, subcommands should be empty (leaf command)
	cmds, err := cmd.Subcommands()
	assert.NoError(t, err)
	assert.Empty(t, cmds)
}

func TestOpenCLICommandDiscovery(t *testing.T) {
	contract := &OpenCLIContract{}
	path := filepath.Join(fixtures, "opencli-tool")
	info, err := os.Lstat(path)
	assert.NoError(t, err)

	d := &discoverer{maxDepth: -1, executor: defaultExecutor}
	cmd, err := contract.BuildCommand(path, fs.FileInfoToDirEntry(info), nil, d)
	assert.NoError(t, err)

	// Summary comes from OpenCLI document
	summary, err := cmd.Summary()
	assert.NoError(t, err)
	assert.Equal(t, "An OpenCLI tool", summary)

	// Subcommands are discovered from OpenCLI output
	cmds, err := cmd.Subcommands()
	assert.NoError(t, err)
	assert.Len(t, cmds, 3)
	assert.Equal(t, "build", cmds[0].Name())
	assert.Equal(t, "mod", cmds[1].Name())
	assert.Equal(t, "hidden-cmd", cmds[2].Name())

	// Nested subcommands
	modCmds, err := cmds[1].Subcommands()
	assert.NoError(t, err)
	assert.Len(t, modCmds, 2)
	assert.Equal(t, "init", modCmds[0].Name())
	assert.Equal(t, "tidy", modCmds[1].Name())

	// Aliases are propagated
	assert.Equal(t, []string{"t"}, modCmds[1].Aliases())

	// DefaultSubcommand
	assert.NotNil(t, cmds[1].DefaultSubcommand())
	assert.Equal(t, "tidy", cmds[1].DefaultSubcommand().Name())
}

func TestParseOpenCLI(t *testing.T) {
	cmd := &executableCommand{path: "/test"}
	out := `{
		"opencli": "0.1-block.1",
		"name": "myapp",
		"info": {"version": "1.0.0"},
		"summary": "My app summary",
		"defaultCommand": "run",
		"commands": [
			{
				"name": "run",
				"summary": "Run the app",
				"aliases": ["r"]
			},
			{
				"name": "build",
				"summary": "Build the app",
				"arguments": [{"name": "target"}],
				"options": [{"name": "--verbose"}],
				"exitCodes": [{"code": 1}]
			}
		]
	}`

	descriptor, err := parseOpenCLI(cmd, out)
	assert.NoError(t, err)
	assert.Equal(t, "myapp", descriptor.Name)
	assert.Equal(t, "My app summary", *descriptor.Summary)
	assert.Equal(t, "run", descriptor.DefaultCommand)
	assert.Len(t, descriptor.Commands, 2)
	assert.Equal(t, "run", descriptor.Commands[0].Name)
	assert.Equal(t, []string{"r"}, descriptor.Commands[0].Aliases)
	assert.Equal(t, "build", descriptor.Commands[1].Name)
}

func TestParseOpenCLIRejectsInvalidJSON(t *testing.T) {
	cmd := &executableCommand{path: "/test"}
	_, err := parseOpenCLI(cmd, "not json")
	assert.Error(t, err)

	var cde CommandDescribeError
	assert.ErrorAs(t, err, &cde)
}

func TestParseOpenCLIRejectsMissingVersion(t *testing.T) {
	cmd := &executableCommand{path: "/test"}
	_, err := parseOpenCLI(cmd, `{"name":"myapp","info":{"version":"1.0.0"}}`)
	assert.Error(t, err)
}

func TestOpenCLIToDescriptorIgnoresExtraFields(t *testing.T) {
	cmd := &executableCommand{path: "/test"}
	out := `{
		"opencli": "0.1-block.1",
		"name": "myapp",
		"info": {"version": "1.0.0"},
		"summary": "A tool",
		"interactive": true,
		"options": [{"name": "--help", "aliases": ["-h"]}],
		"conventions": {"optionSeparator": "="},
		"commands": [
			{
				"name": "sub",
				"hidden": true,
				"examples": ["myapp sub"],
				"metadata": [{"name": "key", "value": "val"}],
				"arguments": [{"name": "file"}],
				"options": [{"name": "--force"}]
			}
		]
	}`

	descriptor, err := parseOpenCLI(cmd, out)
	assert.NoError(t, err)
	assert.Equal(t, "myapp", descriptor.Name)
	assert.Equal(t, "A tool", *descriptor.Summary)
	assert.Len(t, descriptor.Commands, 1)
	assert.Equal(t, "sub", descriptor.Commands[0].Name)
}
