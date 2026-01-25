package exoskeleton

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuggestionsFor(t *testing.T) {
	entrypoint := &Entrypoint{name: "e"}
	spec := &directoryCommand{executableCommand: executableCommand{parent: entrypoint, name: "spec"}}
	echoargs := &executableCommand{parent: spec, name: "echoargs"}
	spec.cmds = Commands{echoargs}
	// Should never be returned because `spec` precedes it.
	overloaded_spec := &executableCommand{parent: entrypoint, name: "spec"}
	entrypoint.cmds = Commands{spec, overloaded_spec}

	scenarios := []struct {
		typedName           string
		expectedSuggestions []Command
	}{
		// When given just part of the name of a command, it suggests the rest
		{"s", []Command{spec, echoargs}},

		// When given a typoed name, it suggest the correct spelling
		{"spec ehcoargs", []Command{echoargs}},
		// ...and this works with colon-separated command names as well
		{"spec:ehcoargs", []Command{echoargs}},

		// When given the name of a command without its namespace, it suggests the namespace
		{"echoargs", []Command{echoargs}},
	}

	for _, s := range scenarios {
		assert.Equal(t, s.expectedSuggestions, entrypoint.suggestionsFor(s.typedName), fmt.Sprintf("SuggestionsFor(\"%s\")", s.typedName))
	}
}
