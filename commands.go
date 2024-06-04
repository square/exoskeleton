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
// with their subcommands. Flatten discards errors.
func (c Commands) Flatten() Commands {
	all := Commands{}
	for _, cmd := range c {
		all = append(all, cmd)
		if m, ok := cmd.(Module); ok {
			subcmds, _ := m.Subcommands()
			all = append(all, subcmds.Flatten()...)
		}
	}
	return all
}
