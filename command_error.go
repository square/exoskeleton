package exoskeleton

// CommandError records an error that occurred with a command's implementation of its interface
type CommandError struct {
	Message string
	Cause   error
	Command Command
}

func (e CommandError) Error() string { return e.Message }
func (e CommandError) Unwrap() error { return e.Cause }
