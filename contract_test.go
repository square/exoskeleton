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
