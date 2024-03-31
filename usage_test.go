package exoskeleton

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsageRelativeTo(t *testing.T) {
	entrypoint := &Entrypoint{name: "e"}
	// └── a
	//     └── b
	//         └── c
	a := &directoryModule{executableCommand: executableCommand{parent: entrypoint, name: "a"}}
	b := &directoryModule{executableCommand: executableCommand{parent: a, name: "b"}}
	c := &executableCommand{parent: b, name: "c"}
	a.cmds = Commands{b}
	b.cmds = Commands{c}
	entrypoint.cmds = Commands{a}

	assert.Equal(t, "e", Usage(entrypoint))

	assert.Equal(t, "e a", Usage(a))
	assert.Equal(t, "e a", UsageRelativeTo(a, nil))
	assert.Equal(t, "a", UsageRelativeTo(a, entrypoint))

	assert.Equal(t, "e a b c", Usage(c))
	assert.Equal(t, "e a b c", UsageRelativeTo(c, nil))
	assert.Equal(t, "a b c", UsageRelativeTo(c, entrypoint))
	assert.Equal(t, "b c", UsageRelativeTo(c, a))
	assert.Equal(t, "c", UsageRelativeTo(c, b))

	// b is not an ancestor of a
	assert.Equal(t, "e a", UsageRelativeTo(a, b))

	// A Command's usage relative to itself is ""
	assert.Equal(t, "", UsageRelativeTo(entrypoint, entrypoint))
}
