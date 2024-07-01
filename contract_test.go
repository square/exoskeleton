package exoskeleton

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessageFromExecutionHandlesErrors(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "edge-cases", "summary-fail")}
	summary, err := getMessageFromExecution(cmd, "summary")

	var ee *exec.ExitError

	assert.Equal(t, "out", summary)
	assert.ErrorAs(t, err, &ee)
	assert.Equal(t, "err\n", string(ee.Stderr))
}

func TestGetMessageFromExecutionTrimsLineBreaks(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "edge-cases", "summary-with-newlines")}
	summary, err := getMessageFromExecution(cmd, "summary")

	assert.NoError(t, err)
	assert.Equal(t, "out", summary)
}

func TestGetMessageFromExecutionWithoutArgs(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "echoargs")}
	output, err := getMessageFromExecution(cmd, "summary")

	assert.NoError(t, err)
	assert.Equal(t, "--summary", output)
}

func TestGetMessageFromExecutionWithArgs(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "echoargs"), args: []string{"a", "b"}}
	output, err := getMessageFromExecution(cmd, "summary")

	assert.NoError(t, err)
	assert.Equal(t, "a\nb\n--summary", output)
}
