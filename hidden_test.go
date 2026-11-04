package exoskeleton

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHiddenDescriptorCommandsStayCallable(t *testing.T) {
	dir := t.TempDir()
	script := `#!/bin/sh
if [ "$1" = "--describe-commands" ]; then
  echo '{"name":"tool","summary":"Example tool","commands":[{"name":"visible","summary":"Visible command"},{"name":"compat","aliases":["old"],"summary":"Compatibility command","hidden":true},{"name":"internal","summary":"Internal commands","hidden":true,"commands":[{"name":"nested","summary":"Nested command"}]}]}'
else
  printf '%s\n' "$@"
fi
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "tool.exoskeleton"), []byte(script), 0o755))
	e, err := New([]string{dir}, WithName("example"))
	require.NoError(t, err)
	tool, _, err := e.Identify([]string{"tool"})
	require.NoError(t, err)
	for _, depth := range []int{0, -1} {
		help, errs := MenuFor(tool, &MenuOptions{Depth: depth})
		require.Empty(t, errs)
		assert.Contains(t, help, "visible")
		for _, name := range []string{"compat", "internal", "nested"} {
			assert.NotContains(t, help, name)
		}
	}
	for _, prefix := range []string{"", "co", "ol"} {
		completions, _, err := tool.Complete(e, []string{prefix}, nil)
		require.NoError(t, err)
		assert.NotContains(t, completions, "compat")
		assert.NotContains(t, completions, "old")
		assert.NotContains(t, completions, "internal")
		if prefix == "" {
			assert.Equal(t, []string{"visible"}, completions)
		}
	}
	for _, name := range []string{"compat", "old"} {
		cmd, rest, err := e.Identify([]string{"tool", name, "argument"})
		require.NoError(t, err)
		assert.Equal(t, "compat", cmd.Name())
		assert.Equal(t, []string{"argument"}, rest)
		executable := cmd.(*executableCommand)
		out, err := executable.output(executable.Command(rest...))
		require.NoError(t, err)
		assert.Equal(t, "compat\nargument\n", string(out))
	}
}
