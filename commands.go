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
// with their subcommands, along with any errors returned by modules'
// Subcommands().
func (c Commands) Flatten() (Commands, []error) {
	all := Commands{}
	errs := []error{}

	for _, cmd := range c {
		all = append(all, cmd)
		if m, ok := cmd.(Module); ok {
			subcmds, err := m.Subcommands()
			if err != nil {
				errs = append(errs, err)
			}
			fcmds, ferrs := subcmds.Flatten()
			all = append(all, fcmds...)
			errs = append(errs, ferrs...)
		}
	}

	return all, errs
}
