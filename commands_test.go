package exoskeleton

import (
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

func TestFlatten(t *testing.T) {
	a := &executableCommand{name: "a"}
	b := &directoryModule{executableCommand: executableCommand{name: "b"}}
	c := &executableCommand{parent: b, name: "c"}
	d := &directoryModule{executableCommand: executableCommand{parent: b, name: "d"}}
	e := &executableCommand{parent: d, name: "e"}
	b.cmds = Commands{c, d}
	d.cmds = Commands{e}

	given := Commands{a, b}
	expected := Commands{a, b, c, d, e}

	assert.Equal(t, expected, given.Flatten())
}
