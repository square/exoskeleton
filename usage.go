package exoskeleton

import "strings"

// Usage returns the usage string for the given command.
// For example, given the Tidy command in the Go CLI, Usage(Tidy) would be 'go mod tidy'.
func Usage(cmd Command) string {
	return UsageRelativeTo(cmd, nil)
}

// UsageRelativeTo returns the usage string for the given command relative to the given command.
// For example, given the Tidy command in the Go CLI ('go mod tidy'), UsageRelativeTo(Tidy, Mod)
// would be 'tidy' and UsageRelativeTo(Tidy, Go) would be 'mod tidy'.
func UsageRelativeTo(cmd Command, relativeTo Command) string {
	return strings.Join(argsRelativeTo(cmd, relativeTo), " ")
}

func argsRelativeTo(cmd Command, relativeTo Command) []string {
	args := []string{}
	for parent := cmd; parent != nil && parent != relativeTo; parent = parent.Parent() {
		args = append([]string{parent.Name()}, args...)
	}
	return args
}
