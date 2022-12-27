package exoskeleton

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscoverIn(t *testing.T) {
	scenarios := []struct {
		maxDepth int
		expected []string
	}{
		{
			0,
			[]string{
				"echoargs",
				"env",
				"exit",
				"hello",
				"suggest",
			},
		},
		{
			1,
			[]string{
				"echoargs",
				"env",
				"exit",
				"hello",
				"nested-1",
				"nested-1 hello",
				"suggest",
				"symlink-test",
				"symlink-test hello-prime",
			},
		},
		{
			2,
			[]string{
				"echoargs",
				"env",
				"exit",
				"hello",
				"nested-1",
				"nested-1 hello",
				"nested-1 nested-2",
				"nested-1 nested-2 hello",
				"suggest",
				"symlink-test",
				"symlink-test hello-prime",
			},
		},
		{
			-1,
			[]string{
				"echoargs",
				"env",
				"exit",
				"hello",
				"nested-1",
				"nested-1 hello",
				"nested-1 nested-2",
				"nested-1 nested-2 hello",
				"nested-1 nested-2 nested-3",
				"nested-1 nested-2 nested-3 hello",
				"nested-1 nested-2 nested-3 nested-4",
				"nested-1 nested-2 nested-3 nested-4 hello",
				"suggest",
				"symlink-test",
				"symlink-test hello-prime",
			},
		},
	}

	for _, s := range scenarios {
		var cmds Commands
		d := discoverer{onError: func(_ error) {}, modulefile: ".exoskeleton", maxDepth: s.maxDepth}
		d.discoverIn(fixtures, nil, &cmds)

		var names []string
		for _, cmd := range cmds.Flatten() {
			names = append(names, Usage(cmd))
		}

		assert.Equal(t, s.expected, names, "When maxDepth=%d", s.maxDepth)
	}
}
func TestFollowingSymlinks(t *testing.T) {
	var cmds Commands
	d := discoverer{onError: func(_ error) {}, modulefile: ".exoskeleton", maxDepth: 1}
	d.discoverIn(fixtures, nil, &cmds)

	cmd := cmds.Find("symlink-test").(Module).Subcommands().Find("hello-prime")

	assert.FileExists(t, cmd.Path())
}
