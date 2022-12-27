package exoskeleton

type Commands []Command

// Find returns the first command with a given name.
func (c Commands) Find(cmdName string) Command {
	for _, cmd := range c {
		if cmd.Name() == cmdName {
			return cmd
		}
	}
	return nil
}

// Flatten returns a list of commands, recursively replacing modules
// with their subcommands.
func (c Commands) Flatten() (all Commands) {
	for _, cmd := range c {
		all = append(all, cmd)
		if m, ok := cmd.(Module); ok {
			all = append(all, m.Subcommands().Flatten()...)
		}
	}
	return
}
