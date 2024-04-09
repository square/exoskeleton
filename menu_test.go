package exoskeleton

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMenuUsage(t *testing.T) {
	entrypoint := &Entrypoint{name: "entrypoint"}
	module := &directoryModule{executableCommand: executableCommand{parent: entrypoint, name: "module"}}
	entrypoint.cmds = Commands{module}

	assert.Equal(t, "entrypoint <command> [<args>]", entrypoint.buildMenu(entrypoint.cmds, entrypoint).Usage)
	assert.Equal(t, "entrypoint module <command> [<args>]", entrypoint.buildMenu(module.cmds, module).Usage)
}

func TestBuildMenuTrailer(t *testing.T) {
	entrypoint := &Entrypoint{name: "entrypoint"}
	module := &directoryModule{executableCommand: executableCommand{parent: entrypoint, name: "module"}}
	entrypoint.cmds = Commands{module}

	assert.Equal(t,
		"Run \033[96mentrypoint help <command>\033[0m to print information on a specific command.",
		entrypoint.buildMenu(entrypoint.cmds, entrypoint).Trailer)

	assert.Equal(t,
		"Run \033[96mentrypoint help module <command>\033[0m to print information on a specific command.",
		entrypoint.buildMenu(module.cmds, module).Trailer)
}

func TestBuildMenuSectionsUncached(t *testing.T) {
	entrypoint, err := New([]string{fixtures})
	entrypoint.cachePath = ""
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,
		`COMMANDS
   echoargs  Echoes the args it received
   env       Prints environment variables
   exit      Exits with the given code
   go:       Provides several commands
   hello     Prints "hello"
   suggest   Suggests arguments`,
		nocolor(entrypoint.buildMenu(entrypoint.cmds, entrypoint).Sections.String()))
}

func TestBuildMenuSectionsReadFromCache(t *testing.T) {
	entrypoint := newWithDefaults("/entrypoint")
	cmd := &executableCommand{parent: entrypoint, path: filepath.Join(fixtures, "echoargs"), name: "echoargs"}
	modTime, err := modTime(cmd)
	if err != nil {
		t.Error(err)
	}
	entrypoint.cmds = Commands{cmd}

	f, err := os.CreateTemp("", "exoskeleton-cache.json")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	f.Write([]byte(fmt.Sprintf(`{"summary":{"entrypoint echoargs":{"modTime":%d,"value":"CACHED SUMMARY"}}}`, modTime)))
	entrypoint.cachePath = f.Name()

	assert.Equal(t,
		"COMMANDS\n   echoargs  CACHED SUMMARY",
		nocolor(entrypoint.buildMenu(entrypoint.cmds, entrypoint).Sections.String()))
}

func TestBuildMenuSectionsWriteToCacheWhenStale(t *testing.T) {
	entrypoint := newWithDefaults("/entrypoint")
	cmd := &executableCommand{parent: entrypoint, path: filepath.Join(fixtures, "echoargs"), name: "echoargs"}
	modTime, err := modTime(cmd)
	if err != nil {
		t.Error(err)
	}
	entrypoint.cmds = Commands{cmd}

	f, err := os.CreateTemp("", "exoskeleton-cache.json")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	f.Write([]byte(fmt.Sprintf(`{"summary":{"entrypoint echoargs":{"modTime":%d,"value":"STALE SUMMARY"}}}`, modTime-1)))
	entrypoint.cachePath = f.Name()

	assert.Equal(t,
		"COMMANDS\n   echoargs  Echoes the args it received",
		nocolor(entrypoint.buildMenu(entrypoint.cmds, entrypoint).Sections.String()))

	f.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)

	assert.Equal(t, fmt.Sprintf(`{"summary":{"entrypoint echoargs":{"modTime":%d,"value":"Echoes the args it received"}}}`, modTime), buf.String())
}

func TestBuildMenuSectionsWriteToCacheWhenMissing(t *testing.T) {
	entrypoint := newWithDefaults("/entrypoint")
	cmd := &executableCommand{parent: entrypoint, path: filepath.Join(fixtures, "echoargs"), name: "echoargs"}
	modTime, err := modTime(cmd)
	if err != nil {
		t.Error(err)
	}
	entrypoint.cmds = Commands{cmd}

	f, err := os.CreateTemp("", "exoskeleton-cache.json")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	f.Write([]byte(`{"summary":{}`))
	entrypoint.cachePath = f.Name()

	assert.Equal(t,
		"COMMANDS\n   echoargs  Echoes the args it received",
		nocolor(entrypoint.buildMenu(entrypoint.cmds, entrypoint).Sections.String()))

	f.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)

	assert.Equal(t, fmt.Sprintf(`{"summary":{"entrypoint echoargs":{"modTime":%d,"value":"Echoes the args it received"}}}`, modTime), buf.String())
}

// TODO: support NOCOLOR and remove this
func nocolor(s string) string {
	re := regexp.MustCompile("\033\\[[;\\d]+m")
	return re.ReplaceAllLiteralString(s, "")
}
