package exoskeleton

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	a := &executable{name: "a"}
	b := &executable{name: "b"}
	// Should never be returned because `a` precedes it.
	overloaded_a := &executable{name: "a"}
	cmds := Commands{a, b, overloaded_a}

	assert.Equal(t, a, cmds.Find("a"))
	assert.Equal(t, b, cmds.Find("b"))
	assert.Equal(t, nil, cmds.Find("c"))
}

func TestFlatten(t *testing.T) {
	a := &executable{name: "a"}
	b := &module{executable: executable{name: "b"}}
	c := &executable{parent: b, name: "c"}
	d := &module{executable: executable{parent: b, name: "d"}}
	e := &executable{parent: d, name: "e"}
	b.cmds = Commands{c, d}
	d.cmds = Commands{e}

	given := Commands{a, b}
	expected := Commands{a, b, c, d, e}

	assert.Equal(t, expected, given.Flatten())
}
