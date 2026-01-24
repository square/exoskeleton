package exoskeleton

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMenuUsage(t *testing.T) {
	entrypoint := &Entrypoint{name: "entrypoint"}
	module := &directoryModule{executableCommand: executableCommand{parent: entrypoint, name: "module", cache: nullCache{}}}
	entrypoint.cmds = Commands{module}

	menu, _ := buildMenu(entrypoint, &MenuOptions{})
	assert.Equal(t, "entrypoint <command> [<args>]", menu.Usage)

	menu, _ = buildMenu(module, &MenuOptions{})
	assert.Equal(t, "entrypoint module <command> [<args>]", menu.Usage)
}

func TestMenuForTrailer(t *testing.T) {
	entrypoint := &Entrypoint{name: "entrypoint"}
	module := &directoryModule{executableCommand: executableCommand{parent: entrypoint, name: "module", cache: nullCache{}}}
	entrypoint.cmds = Commands{module}

	menu, _ := MenuFor(entrypoint, &MenuOptions{})
	assert.Contains(t, menu, "Run \033[96mentrypoint help <command>\033[0m to print information on a specific command.")

	menu, _ = MenuFor(module, &MenuOptions{})
	assert.Contains(t, menu, "Run \033[96mentrypoint help module <command>\033[0m to print information on a specific command.")
}

func TestMenuForSections(t *testing.T) {
	entrypoint, err := New([]string{fixtures})
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
		menu, errs := MenuFor(entrypoint, &MenuOptions{Depth: s.depth})
		assert.Empty(t, errs)
		assert.Equal(t, s.expected, sections(nocolor(menu)), "Given depth=%d", s.depth)
	}
}

// TODO: support NOCOLOR and remove this
func nocolor(s string) string {
	re := regexp.MustCompile("\033\\[[;\\d]+m")
	return re.ReplaceAllLiteralString(s, "")
}

func sections(s string) string {
	lines := strings.SplitAfter(s, "\n")
	return strings.TrimRight(strings.Join(lines[3:(len(lines)-1)], ""), "\n")
}
