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
	return c.Expand()
}

// Expand returns a list of commands, recursively replacing modules
// with their subcommands up to a given depth, along with any errors
// returned by modules' Subcommands().
func (c Commands) Expand(fops ...ExpandOption) (Commands, []error) {
	o := &expandOptions{
		depth:                  -1,
		includeExpandedModules: true,
	}
	for _, fop := range fops {
		fop(o)
	}
	return expand(c, o.depth, o.includeExpandedModules)
}

type ExpandOption func(*expandOptions)

type expandOptions struct {
	depth                  int
	includeExpandedModules bool
}

func WithDepth(d int) ExpandOption {
	return func(o *expandOptions) { o.depth = d }
}

func WithoutExpandedModules() ExpandOption {
	return func(o *expandOptions) { o.includeExpandedModules = false }
}

func expand(c Commands, depth int, includeExpandedModules bool) (Commands, []error) {
	all := Commands{}
	errs := []error{}

	for _, cmd := range c {
		if m, ok := cmd.(Module); ok && depth != 0 {
			if includeExpandedModules {
				all = append(all, cmd)
			}

			subcmds, err := m.Subcommands()
			if err != nil {
				errs = append(errs, err)
			}
			fcmds, ferrs := expand(subcmds, depth-1, includeExpandedModules)
			all = append(all, fcmds...)
			errs = append(errs, ferrs...)
		} else {
			all = append(all, cmd)
		}
	}

	return all, errs
}
