package exoskeleton

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNull(t *testing.T) {
	assert.True(t, IsNull(nullCommand{}))
	assert.False(t, IsNull(&executableCommand{}))
	assert.False(t, IsNull(&builtinCommand{}))
}
