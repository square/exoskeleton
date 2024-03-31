package exoskeleton

import (
	"fmt"
	"testing"

	"github.com/square/exoskeleton/pkg/shellcomp"
	"github.com/stretchr/testify/assert"
)

func TestCompletionsFor(t *testing.T) {
	// all
	// ├── echo
	// └── s-foils
	//     ├── lock
	//     └── unlock
	entrypoint := &Entrypoint{}
	echo := &builtinCommand{parent: entrypoint, name: "echo", complete: echoArgs}
	sfoils := &directoryModule{executableCommand: executableCommand{parent: entrypoint, name: "s-foils"}}
	lock := &executableCommand{parent: sfoils, name: "lock"}
	unlock := &executableCommand{parent: sfoils, name: "unlock"}
	sfoils.cmds = Commands{lock, unlock}
	entrypoint.cmds = Commands{echo, sfoils}

	scenarios := []struct {
		args                []string
		completeArgs        bool
		expectedCompletions []string
	}{
		// Should suggest all commands when the user hasn't typed anything yet
		{[]string{""}, true, []string{"echo", "s-foils"}},
		{[]string{"s-foils", ""}, true, []string{"lock", "unlock"}},

		// Should suggest only commands that start with whatever the user typed
		{[]string{"e"}, true, []string{"echo"}},
		{[]string{"s"}, true, []string{"s-foils"}},
		{[]string{"s-foils", "unloc"}, true, []string{"unlock"}},

		// Should suggest what the user typed if they typed the exact name of a command
		{[]string{"s-foils", "unlock"}, true, []string{"unlock"}},

		// Should delegate completions for arguments and flags to commands themselves
		{[]string{"echo", "a", "b"}, true, []string{"a", "b"}},
		{[]string{"echo", "--flag"}, true, []string{"--flag"}},

		// Should return no suggestions if argument-completion is not requested
		{[]string{"echo", "a", "b"}, false, nil},
		{[]string{"echo", "--flag"}, false, nil},
	}

	for _, s := range scenarios {
		actualCompletions, _, _ := entrypoint.completionsFor(s.args, nil, s.completeArgs)
		assert.Equal(t, s.expectedCompletions, actualCompletions, fmt.Sprintf("completionsFor(\"%v\", %t)", s.args, s.completeArgs))
	}
}

func echoArgs(_ *Entrypoint, args, _ []string) ([]string, shellcomp.Directive, error) {
	return args, shellcomp.DirectiveDefault, nil
}
