package exoskeleton

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpForWithMagicComment(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "edge-cases", "help-from-magic-comments")}
	entrypoint := &Entrypoint{cmds: Commands{cmd}}
	help, err := entrypoint.helpFor(cmd, nil)

	assert.NoError(t, err)
	assert.Equal(t, "USAGE: help-from-magic-comments", help)
}

func TestHelpForWithExecution(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "edge-cases", "help-from-execution")}
	entrypoint := &Entrypoint{cmds: Commands{cmd}}
	help, err := entrypoint.helpFor(cmd, nil)

	assert.NoError(t, err)
	assert.Equal(t, "USAGE: help from execution", help)
}

func TestHelpForWithExecutionPassesArgsThrough(t *testing.T) {
	cmd := &executableCommand{path: filepath.Join(fixtures, "echoargs"), args: []string{"a"}}
	entrypoint := &Entrypoint{cmds: Commands{cmd}}
	help, err := entrypoint.helpFor(cmd, []string{"b"})

	assert.NoError(t, err)
	assert.Equal(t, "a\nb\n--help", help)
}
