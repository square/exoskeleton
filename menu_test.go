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

	menu, _ := buildMenu(entrypoint, &buildMenuOptions{})
	assert.Equal(t, "entrypoint <command> [<args>]", menu.Usage)

	menu, _ = buildMenu(module, &buildMenuOptions{})
	assert.Equal(t, "entrypoint module <command> [<args>]", menu.Usage)
}

func TestBuildMenuTrailer(t *testing.T) {
	entrypoint := &Entrypoint{name: "entrypoint"}
	module := &directoryModule{executableCommand: executableCommand{parent: entrypoint, name: "module"}}
	entrypoint.cmds = Commands{module}

	menu, _ := buildMenu(entrypoint, &buildMenuOptions{})
	assert.Contains(t, menu.String(), "Run \033[96mentrypoint help <command>\033[0m to print information on a specific command.")

	menu, _ = buildMenu(module, &buildMenuOptions{})
	assert.Contains(t, menu.String(), "Run \033[96mentrypoint help module <command>\033[0m to print information on a specific command.")
}

func TestBuildMenuSectionsUncached(t *testing.T) {
	entrypoint, err := New([]string{fixtures})
	entrypoint.cachePath = ""
	if err != nil {
		t.Error(err)
	}

	scenarios := []struct {
		depth    int
		expected string
	}{
		{
			0,
			`COMMANDS
   echoargs  Echoes the args it received
   env       Prints environment variables
   exit      Exits with the given code
   go:       Provides several commands
   hello     Prints "hello"
   suggest   Suggests arguments`,
		},
		{
			1,
			`COMMANDS
   echoargs  Echoes the args it received
   env       Prints environment variables
   exit      Exits with the given code
   go build  compile packages and dependencies
   go mod:   module maintenance
   hello     Prints "hello"
   suggest   Suggests arguments`,
		},
		{
			2,
			`COMMANDS
   echoargs     Echoes the args it received
   env          Prints environment variables
   exit         Exits with the given code
   go build     compile packages and dependencies
   go mod init  initialize new module in current directory
   go mod tidy  add missing and remove unused modules
   hello        Prints "hello"
   suggest      Suggests arguments`,
		},
		{
			-1,
			`COMMANDS
   echoargs     Echoes the args it received
   env          Prints environment variables
   exit         Exits with the given code
   go build     compile packages and dependencies
   go mod init  initialize new module in current directory
   go mod tidy  add missing and remove unused modules
   hello        Prints "hello"
   suggest      Suggests arguments`,
		},
	}

	for _, s := range scenarios {
		menu, errs := buildMenu(entrypoint, &buildMenuOptions{Depth: s.depth})
		assert.Empty(t, errs)
		assert.Equal(t, s.expected, nocolor(menu.Sections.String()), "Given depth=%d", s.depth)
	}
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
	cache := &summaryCache{Path: f.Name(), onError: entrypoint.onError}

	menu, errs := buildMenu(entrypoint, &buildMenuOptions{SummaryFor: cache.Read})
	assert.Empty(t, errs)
	assert.Equal(t, "COMMANDS\n   echoargs  CACHED SUMMARY", nocolor(menu.Sections.String()))
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
	cache := &summaryCache{Path: f.Name(), onError: entrypoint.onError}

	menu, errs := buildMenu(entrypoint, &buildMenuOptions{SummaryFor: cache.Read})
	assert.Empty(t, errs)
	assert.Equal(t, "COMMANDS\n   echoargs  Echoes the args it received", nocolor(menu.Sections.String()))

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
	cache := &summaryCache{Path: f.Name(), onError: entrypoint.onError}

	menu, errs := buildMenu(entrypoint, &buildMenuOptions{SummaryFor: cache.Read})
	assert.Empty(t, errs)
	assert.Equal(t, "COMMANDS\n   echoargs  Echoes the args it received", nocolor(menu.Sections.String()))

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
