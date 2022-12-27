package exoskeleton

import "fmt"

// SymlinkError records an error that occurred while following a symlink.
type SymlinkError struct {
	Cause error
	Path  string
}

func (e SymlinkError) Error() string { return fmt.Sprintf("broken symlink %s: %s", e.Path, e.Cause) }
func (e SymlinkError) Unwrap() error { return e.Cause }
