package shellcomp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	scenarios := []struct {
		completions []string
		directive   Directive
		expected    string
	}{
		{[]string{}, DirectiveNoFileComp, ":4\n"},
		{[]string{}, DirectiveFilterDirs | DirectiveNoSpace, ":18\n"},

		{[]string{"a", "b", "c"}, DirectiveDefault, "a\nb\nc\n:0\n"},
	}

	for _, s := range scenarios {
		b := Marshal(s.completions, s.directive, false)
		assert.Equal(t, s.expected, string(b), fmt.Sprintf("Marshal(%v, %d)", s.completions, s.directive))
	}
}
