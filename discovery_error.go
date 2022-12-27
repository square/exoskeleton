package exoskeleton

import "fmt"

// DiscoveryError records an error that occurred while discovering commands in a directory.
type DiscoveryError struct {
	Cause error
	Path  string
}

func (e DiscoveryError) Error() string {
	return fmt.Sprintf("error discovering commands in %s: %s", e.Path, e.Cause)
}
func (e DiscoveryError) Unwrap() error { return e.Cause }
