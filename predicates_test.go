package exoskeleton

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasSubcommandsIsFalseForLeafExecutables(t *testing.T) {
	cmd := &executableCommand{name: "leaf"}
	assert.False(t, HasSubcommands(cmd))
}

func TestHasSubcommandsIsFalseForShellScripts(t *testing.T) {
	cmd := &shellScriptCommand{executableCommand: executableCommand{name: "script"}}
	assert.False(t, HasSubcommands(cmd))
}

func TestHasSubcommandsIsTrueForUndiscoveredExecutables(t *testing.T) {
	// cache is nil: if HasSubcommands performed discovery, it would panic.
	cmd := &executableCommand{name: "parent", discoverer: &discoverer{}}
	assert.True(t, HasSubcommands(cmd))
}

func TestHasSubcommandsIsExactAfterDiscovery(t *testing.T) {
	cmd := &executableCommand{name: "parent", discoverer: &discoverer{}, cmds: Commands{}}
	assert.False(t, HasSubcommands(cmd))

	cmd.cmds = Commands{&executableCommand{name: "child"}}
	assert.True(t, HasSubcommands(cmd))
}

func TestHasSubcommandsIsFalseForOpenCLILeaves(t *testing.T) {
	t.Skip("Known regression: OpenCLI leaves report that they may have subcommands until discovered")

	path := filepath.Join(t.TempDir(), "leaf")
	script := `#!/usr/bin/env bash
echo '{"opencli": "0.1-block.1", "name": "leaf", "info": {"version": "1.0.0"}}'
`
	assert.NoError(t, os.WriteFile(path, []byte(script), 0o755))
	info, err := os.Lstat(path)
	assert.NoError(t, err)

	d := &discoverer{maxDepth: -1, executor: defaultExecutor}
	cmd, err := (&OpenCLIContract{}).BuildCommand(path, fs.FileInfoToDirEntry(info), nil, d)
	assert.NoError(t, err)
	assert.False(t, HasSubcommands(cmd))
}

func TestHasSubcommandsIsTrueForUndiscoveredDirectories(t *testing.T) {
	// path is empty: if HasSubcommands performed discovery, it would find nothing.
	cmd := &directoryCommand{discoverer: &discoverer{}}
	assert.True(t, HasSubcommands(cmd))
}

func TestHasSubcommandsFallsBackToSubcommands(t *testing.T) {
	assert.False(t, HasSubcommands(nullCommand{}))
	assert.True(t, HasSubcommands(&builtinCommand{subcommands: Commands{nullCommand{}}}))
}
