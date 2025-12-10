package exoskeleton

import "sync"

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
	type result struct {
		commands Commands
		errors   []error
	}
	results := make([]result, len(c))
	var wg sync.WaitGroup

	for i, cmd := range c {
		// If this is a module, recursively flatten its subcommands...
		if m, ok := cmd.(Module); ok && depth != 0 {
			// ...and do that in parallel because resolving subcommands may involve
			// executing them with `--describe-commands` and subcommands that accidentally
			// introduce latency impact the experience.
			wg.Add(1)
			go func(idx int, mod Module) {
				defer wg.Done()
				cmds := Commands{}
				errs := []error{}

				if includeExpandedModules {
					cmds = append(cmds, mod)
				}

				subcmds, err := mod.Subcommands()
				if err != nil {
					errs = append(errs, err)
				}

				fcmds, ferrs := expand(subcmds, depth-1, includeExpandedModules)
				cmds = append(cmds, fcmds...)
				errs = append(errs, ferrs...)

				// Store result at the correct index to maintain order
				results[idx] = result{commands: cmds, errors: errs}
			}(i, m)

		} else {
			results[i] = result{commands: Commands{cmd}, errors: []error{}}
		}
	}

	wg.Wait()

	all := Commands{}
	errs := []error{}
	for _, r := range results {
		all = append(all, r.commands...)
		errs = append(errs, r.errors...)
	}

	return all, errs
}
