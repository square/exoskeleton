package exoskeleton

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscoverInFindsNothing(t *testing.T) {
	entrypoint := &Entrypoint{}
	cmds := entrypoint.discoverIn([]string{filepath.Join(fixtures, "nope")})
	assert.Equal(t, Commands{}, cmds)
}

func TestDiscoverInWithMaxDepth(t *testing.T) {
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
				"go",
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
				"go",
				"go build",
				"go mod",
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
				"go",
				"go build",
				"go mod",
				"go mod init",
				"go mod tidy",
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
				"go",
				"go build",
				"go mod",
				"go mod init",
				"go mod tidy",
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
		d := discoverer{entrypoint: &Entrypoint{moduleMetadataFilename: ".exoskeleton", maxDepth: s.maxDepth}}
		d.discoverIn(fixtures, nil, &cmds)

		var names []string
		for _, cmd := range cmds.Flatten() {
			names = append(names, Usage(cmd))
		}

		assert.Equal(t, s.expected, names, "When maxDepth=%d", s.maxDepth)
	}
}

func TestDiscovererBuildsCommand(t *testing.T) {
	var parent Module = &builtinModule{}
	d := discoverer{entrypoint: &Entrypoint{moduleMetadataFilename: ".exoskeleton", maxDepth: 2}}

	scenarios := []struct {
		executable string
		expected   Command
	}{
		{
			"echoargs",
			&executableCommand{
				entrypoint:   d.entrypoint,
				parent:       parent,
				name:         "echoargs",
				path:         filepath.Join(fixtures, "echoargs"),
				discoveredIn: fixtures,
			},
		},
		{
			"nested-1",
			&directoryModule{
				executableCommand: executableCommand{
					entrypoint:   d.entrypoint,
					parent:       parent,
					name:         "nested-1",
					path:         filepath.Join(fixtures, "nested-1", ".exoskeleton"),
					discoveredIn: fixtures,
				},
				discoverer: discoverer{
					entrypoint: d.entrypoint,
					depth:      d.depth + 1,
				},
			},
		},
		{
			"go.exoskeleton",
			&executableModule{
				executableCommand: executableCommand{
					entrypoint:   d.entrypoint,
					parent:       parent,
					name:         "go",
					path:         filepath.Join(fixtures, "go.exoskeleton"),
					discoveredIn: fixtures,
				},
				discoverer: discoverer{
					entrypoint: d.entrypoint,
					depth:      d.depth + 1,
				},
			},
		},
	}

	for _, s := range scenarios {
		path := filepath.Join(fixtures, s.executable)
		info, err := os.Lstat(path)
		assert.NoErrorf(t, err, "Given executable=%s", s.executable)
		entry := fs.FileInfoToDirEntry(info)
		cmd, err := d.buildCommand(fixtures, parent, entry)
		assert.NoErrorf(t, err, "Given executable=%s", s.executable)

		assert.Equalf(t, s.expected, cmd, "Given executable=%s", s.executable)
	}
}

func TestFollowingSymlinks(t *testing.T) {
	var cmds Commands
	d := discoverer{entrypoint: &Entrypoint{moduleMetadataFilename: ".exoskeleton", maxDepth: 1}}
	d.discoverIn(fixtures, nil, &cmds)

	cmd := cmds.Find("symlink-test").(Module).Subcommands().Find("hello-prime")

	assert.FileExists(t, cmd.Path())
}
