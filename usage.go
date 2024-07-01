package exoskeleton

import "strings"

// Usage returns the usage string for the given command.
// For example, given the Tidy command in the Go CLI, Usage(Tidy) would be 'go mod tidy'.
func Usage(cmd Command) string {
	return UsageRelativeTo(cmd, nil)
}

// UsageRelativeTo returns the usage string for the given command relative to the given module.
// For example, given the Tidy command in the Go CLI ('go mod tidy'), UsageRelativeTo(Tidy, Mod)
// would be 'tidy' and UsageRelativeTo(Tidy, Go) would be 'mod tidy'.
func UsageRelativeTo(cmd Command, m Module) string {
	return strings.Join(argsRelativeTo(cmd, m), " ")
}

func argsRelativeTo(cmd Command, m Module) []string {
	args := []string{}
	for parent := cmd; parent != nil && parent != m; parent = parent.Parent() {
		args = append([]string{parent.Name()}, args...)
	}
	return args
}
