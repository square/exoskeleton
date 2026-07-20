package exoskeleton

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// -s/--follow-symlinks are which's own flags. They must be consumed by which
// and not forwarded to Identify -- regardless of whether they appear before or
// after the command name -- since a trailing flag would make Identify resolve
// the command's default subcommand instead of the command itself.
func TestWhichConsumesFollowSymlinksFlag(t *testing.T) {
	scenarios := []struct {
		args               []string
		wantIdentifyArgs   []string
		wantFollowSymlinks bool
	}{
		{[]string{"tool"}, []string{"tool"}, false},
		{[]string{"tool", "-s"}, []string{"tool"}, true},
		{[]string{"tool", "--follow-symlinks"}, []string{"tool"}, true},
		{[]string{"-s", "tool"}, []string{"tool"}, true},
		{[]string{"--follow-symlinks", "tool"}, []string{"tool"}, true},

		// Everything after `--` is left for the command, flag or not.
		{[]string{"tool", "--", "-s"}, []string{"tool", "--", "-s"}, false},
	}

	for _, s := range scenarios {
		identifyArgs, followSymlinks := splitWhichArgs(s.args)
		assert.Equal(t, s.wantIdentifyArgs, identifyArgs, "which %v", s.args)
		assert.Equal(t, s.wantFollowSymlinks, followSymlinks, "which %v", s.args)
	}
}
