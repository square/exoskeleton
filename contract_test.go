package exoskeleton

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessageFromExecutionHandlesErrors(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "edge-cases", "summary-fail"), executor: defaultExecutor}
	summary, err := getMessageFromExecution(cmd, "summary")

	var ee *exec.ExitError

	assert.Equal(t, "out", summary)
	assert.ErrorAs(t, err, &ee)
	assert.Equal(t, "err\n", string(ee.Stderr))
}

func TestGetMessageFromExecutionTrimsLineBreaks(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "edge-cases", "summary-with-newlines"), executor: defaultExecutor}
	summary, err := getMessageFromExecution(cmd, "summary")

	assert.NoError(t, err)
	assert.Equal(t, "out", summary)
}

func TestGetMessageFromExecutionWithoutArgs(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "echoargs"), executor: defaultExecutor}
	output, err := getMessageFromExecution(cmd, "summary")

	assert.NoError(t, err)
	assert.Equal(t, "--summary", output)
}

func TestGetMessageFromExecutionWithArgs(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "echoargs"), args: []string{"a", "b"}, executor: defaultExecutor}
	output, err := getMessageFromExecution(cmd, "summary")

	assert.NoError(t, err)
	assert.Equal(t, "a\nb\n--summary", output)
}

func TestParseDescribeCommandsWithAliases(t *testing.T) {
	cmd := &executableCommand{path: "/test"}
	out := `{
		"name": "root",
		"commands": [
			{"name": "remove", "aliases": ["rm", "del"], "summary": "Remove something"},
			{"name": "add", "summary": "Add something"}
		]
	}`

	descriptor, err := parseDescribeCommands(cmd, out)
	assert.NoError(t, err)
	assert.Equal(t, "root", descriptor.Name)
	assert.Len(t, descriptor.Commands, 2)
	assert.Equal(t, []string{"rm", "del"}, descriptor.Commands[0].Aliases)
	assert.Nil(t, descriptor.Commands[1].Aliases)
}

func TestToCommandsPropagatesAliases(t *testing.T) {
	parent := &executableCommand{
		path:     "/test",
		executor: defaultExecutor,
		cache:    nullCache{},
	}

	summary := "Remove something"
	descriptors := []*commandDescriptor{
		{Name: "remove", Aliases: []string{"rm"}, Summary: &summary},
		{Name: "add"},
	}

	cmds := toCommands(parent, descriptors, nil, &discoverer{maxDepth: 0})

	assert.Len(t, cmds, 2)

	removeCmd := cmds[0].(*executableCommand)
	assert.Equal(t, "remove", removeCmd.name)
	assert.Equal(t, []string{"rm"}, removeCmd.aliases)

	addCmd := cmds[1].(*executableCommand)
	assert.Equal(t, "add", addCmd.name)
	assert.Nil(t, addCmd.aliases)
}

func TestDefaultCommandFromDescriptor(t *testing.T) {
	parent := &executableCommand{path: "/test", executor: defaultExecutor, cache: nullCache{}}
	out := `{
		"name": "root",
		"defaultCommand": "b",
		"commands": [
			{"name": "a"},
			{"name": "b"}
		]
	}`

	descriptor, err := parseDescribeCommands(parent, out)
	assert.NoError(t, err)
	assert.Equal(t, "b", descriptor.DefaultCommand)

	cmds := toCommands(parent, descriptor.Commands, nil, &discoverer{maxDepth: 0})
	parent.cmds = cmds
	parent.defaultSubcommand = descriptor.DefaultCommand
	parent.discoverer = &discoverer{maxDepth: 0}

	assert.Equal(t, "b", parent.DefaultSubcommand().Name())
}

// TestWithContractsOption verifies that WithContracts replaces the contracts.
func TestWithContractsOption(t *testing.T) {
	// Create a custom contract that only matches files named "custom"
	customContract := &testContract{
		name: "custom-only",
		matches: func(path string, info fs.DirEntry) bool {
			return filepath.Base(path) == "custom"
		},
	}

	e := &Entrypoint{
		contracts: defaultContracts(),
	}

	// Apply WithContracts
	WithContracts(customContract).Apply(e)

	// Verify contracts were replaced
	if len(e.contracts) != 1 {
		t.Errorf("expected 1 contract, got %d", len(e.contracts))
	}
}

// TestContractOrder verifies that contracts are tried in order.
func TestContractOrder(t *testing.T) {
	matchAllCalled := false
	matchNoneCalled := false

	// First contract matches everything
	matchAll := &testContract{
		name: "match-all",
		matches: func(path string, info fs.DirEntry) bool {
			matchAllCalled = true
			return true
		},
	}

	// Second contract should never be reached
	matchNone := &testContract{
		name: "match-none",
		matches: func(path string, info fs.DirEntry) bool {
			matchNoneCalled = true
			return false
		},
	}

	// Build a fake discoverer
	d := &discoverer{
		maxDepth:  -1,
		contracts: []Contract{matchAll, matchNone},
	}

	// Try to build a command using a real file
	path := filepath.Join(fixtures, "hello")
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("Failed to stat test file: %v", err)
	}
	entry := fs.FileInfoToDirEntry(info)

	_, err = d.buildCommand(fixtures, nil, entry)
	if err != nil {
		t.Fatalf("buildCommand failed: %v", err)
	}

	// Verify first contract was called
	if !matchAllCalled {
		t.Error("first contract (match-all) was not called")
	}

	// Verify second contract was NOT called (first one matched)
	if matchNoneCalled {
		t.Error("second contract (match-none) should not be tried when first matches")
	}
}

// testContract is a simple contract for testing.
type testContract struct {
	name    string
	matches func(path string, info fs.DirEntry) bool
}

func (c *testContract) BuildCommand(path string, info fs.DirEntry, parent Command, d DiscoveryContext) (Command, error) {
	if c.matches != nil && !c.matches(path, info) {
		return nil, ErrNotApplicable
	}
	// Return nil command for testing purposes
	return nil, nil
}

func defaultContracts() []Contract {
	return []Contract{
		&DirectoryContract{MetadataFilename: ".exoskeleton"},
		&ExecutableContract{},
		&ShellScriptContract{},
		&StandaloneExecutableContract{},
	}
}
