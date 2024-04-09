package exoskeleton

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessageFromExecution(t *testing.T) {
	path := filepath.Join(fixtures, "edge-cases", "summary-fail")
	summary, err := getMessageFromExecution(path, "summary")

	var ee *exec.ExitError

	assert.Equal(t, "out", summary)
	assert.ErrorAs(t, err, &ee)
	assert.Equal(t, "err\n", string(ee.Stderr))
}
