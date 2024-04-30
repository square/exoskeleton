package exoskeleton

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeArgs(t *testing.T) {
	scenarios := []struct {
		rawArgs      []string
		expectedArgs []string
		memo         string
	}{
		{[]string{}, []string{}, "Does nothing to empty arguments"},
		{[]string{"ssh"}, []string{"ssh"}, "Does nothing to plain arguments"},
		{[]string{"ssh", "--help"}, []string{"help", "ssh"}, "Normalizes `xyz --help` to `help xyz`"},
		{[]string{"ssh", "-h"}, []string{"help", "ssh"}, "Normalizes `xyz -h` to `help xyz`"},
		{[]string{"ssh", "--", "--help"}, []string{"ssh", "--", "--help"}, "Leaves `--help` alone after `--`"},
	}

	for _, s := range scenarios {
		assert.Equal(t, s.expectedArgs, normalizeArgs(s.rawArgs), s.memo)
	}
}

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
	d := &executableModule{executableCommand: executableCommand{parent: b, name: "d"}}
	e := &executableCommand{parent: d, args: []string{"e"}, name: "e"}
	b.cmds = Commands{c, d}
	d.cmds = Commands{e}

	// Should never be returned because `a` precedes it.
	overloaded_a := &executableCommand{name: "a"}

	help := &builtinCommand{definition: &EmbeddedCommand{Name: `help`}}
	entrypoint := &Entrypoint{cmds: Commands{help, a, b, overloaded_a}}

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
	}

	for _, s := range scenarios {
		cmd, rest := entrypoint.Identify(s.args)
		assert.Equal(t, s.expectedCmd, cmd, fmt.Sprintf("Identify(%v)", s.args))
		assert.Equal(t, s.expectedArgs, rest, fmt.Sprintf("Identify(%v)", s.args))
	}
}
