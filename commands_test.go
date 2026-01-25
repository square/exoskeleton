package exoskeleton

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	a := &executableCommand{name: "a"}
	b := &executableCommand{name: "b"}
	// Should never be returned because `a` precedes it.
	overloaded_a := &executableCommand{name: "a"}
	cmds := Commands{a, b, overloaded_a}

	assert.Equal(t, a, cmds.Find("a"))
	assert.Equal(t, b, cmds.Find("b"))
	assert.Equal(t, nil, cmds.Find("c"))
}

func TestExpand(t *testing.T) {
	a := &executableCommand{name: "a"}
	b := &directoryCommand{executableCommand: executableCommand{name: "b"}}
	c := &executableCommand{parent: b, name: "c"}
	d := &directoryCommand{executableCommand: executableCommand{parent: b, name: "d"}}
	e := &executableCommand{parent: d, name: "e"}
	b.cmds = Commands{c, d}
	d.cmds = Commands{e}
	given := Commands{a, b}

	scenarios :=
		[]struct {
			ops      []ExpandOption
			expected Commands
		}{
			{[]ExpandOption{WithDepth(0)}, Commands{a, b}},
			{[]ExpandOption{WithDepth(1)}, Commands{a, b, c, d}},
			{[]ExpandOption{WithDepth(2)}, Commands{a, b, c, d, e}},
			{[]ExpandOption{WithDepth(-1)}, Commands{a, b, c, d, e}},
			{[]ExpandOption{WithoutExpandedModules()}, Commands{a, c, e}},
		}

	for i, s := range scenarios {
		cmds, errs := given.Expand(s.ops...)
		assert.Empty(t, errs)
		assert.Equal(t, namesOf(s.expected), namesOf(cmds), "Expand[%d]", i)
	}
}

func namesOf(cmds Commands) string {
	var result []string
	for _, cmd := range cmds {
		result = append(result, cmd.Name())
	}
	return strings.Join(result, "\n")
}
