package exoskeleton

import "fmt"

// CacheError records an error that prevented the menu from leverage its cache
type CacheError struct {
	Message string
	Cause   error
}

func (e CacheError) Error() string { return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error()) }
func (e CacheError) Unwrap() error { return e.Cause }
