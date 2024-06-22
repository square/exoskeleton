package exoskeleton

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentify(t *testing.T) {
	// all
	// ├── a
	// └── b
	//     ├── c
	//     └── d
	//         └── e
	a := &executableCommand{name: "a"}
	b := &directoryModule{executableCommand: executableCommand{name: "b"}}
	c := &executableCommand{parent: b, name: "c"}
	d := &directoryModule{executableCommand: executableCommand{parent: b, name: "d"}}
	e := &executableCommand{parent: d, name: "e"}
	b.cmds = Commands{c, d}
	d.cmds = Commands{e}

	// Should never be returned because `a` precedes it.
	overloaded_a := &executableCommand{name: "a"}

	help := &builtinCommand{definition: &EmbeddedCommand{Name: `help`}}
	complete := &builtinCommand{definition: &EmbeddedCommand{Name: `complete`}}
	entrypoint := &Entrypoint{cmds: Commands{help, complete, a, b, overloaded_a}}

	scenarios := []struct {
		args         []string
		expectedCmd  Command
		expectedArgs []string
	}{
		{[]string{}, entrypoint, []string{}},
		{[]string{"x"}, nullCommand{entrypoint, "x"}, []string{}},

		// Can invoke a command (with or without args)
		{[]string{"a"}, a, []string{}},
		{[]string{"a", "arg", "--flag"}, a, []string{"arg", "--flag"}},

		// Can invoke a subcommand the explicit way (with or without args)
		{[]string{"b:c"}, c, []string{}},
		{[]string{"b:c", "arg", "--flag"}, c, []string{"arg", "--flag"}},
		{[]string{"b:d:e"}, e, []string{}},
		{[]string{"b:d:e", "arg", "--flag"}, e, []string{"arg", "--flag"}},

		// Can invoke a subcommand (with or without args)
		{[]string{"b", "c"}, c, []string{}},
		{[]string{"b", "c", "arg", "--flag"}, c, []string{"arg", "--flag"}},
		{[]string{"b", "d", "e"}, e, []string{}},
		{[]string{"b", "d", "e", "arg", "--flag"}, e, []string{"arg", "--flag"}},

		// Can invoke a module with flags but not with positional args
		{[]string{"b"}, b, []string{}},
		{[]string{"b", "--flag"}, b, []string{"--flag"}},
		{[]string{"b", "x", "--flag"}, nullCommand{b, "x"}, []string{"--flag"}},
		{[]string{"b", "d"}, d, []string{}},
		{[]string{"b", "d", "--flag"}, d, []string{"--flag"}},
		{[]string{"b", "d", "x", "--flag"}, nullCommand{d, "x"}, []string{"--flag"}},

		// Can invoke a module with a trailing colon
		{[]string{"b:"}, b, []string{}},

		// Identifies `--complete` as `complete`
		{[]string{"--complete"}, complete, []string{}},

		// Normalizes `CMD --help` to `help CMD`
		{[]string{"--help"}, help, []string{}},
		{[]string{"a", "--help"}, help, []string{"a"}},
		{[]string{"b", "c", "--help"}, help, []string{"b", "c"}},
		{[]string{"b:c", "--help"}, help, []string{"b:c"}},

		// Normalizes `CMD -h` to `help CMD`
		{[]string{"-h"}, help, []string{}},
		{[]string{"a", "-h"}, help, []string{"a"}},
		{[]string{"b", "c", "-h"}, help, []string{"b", "c"}},
		{[]string{"b:c", "-h"}, help, []string{"b:c"}},

		// Does not normalize `--help`/`-h` after `--`
		{[]string{"a", "--", "--help"}, a, []string{"--", "--help"}},
		{[]string{"a", "--", "-h"}, a, []string{"--", "-h"}},
	}

	for _, s := range scenarios {
		cmd, rest, err := entrypoint.Identify(s.args)
		assert.NoError(t, err, fmt.Sprintf("Identify(%v)", s.args))
		assert.Equal(t, s.expectedCmd, cmd, fmt.Sprintf("Identify(%v)", s.args))
		assert.Equal(t, s.expectedArgs, rest, fmt.Sprintf("Identify(%v)", s.args))
	}
}
