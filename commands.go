package exoskeleton

import (
	"sync"
)

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
	return parallelMap(c, func(cmd Command) ([]Command, []error) {
		// If this is a module, recursively flatten its subcommands...
		if m, ok := cmd.(Module); ok && depth != 0 {
			cmds := []Command{}
			errs := []error{}

			if includeExpandedModules {
				cmds = append(cmds, m)
			}

			subcmds, err := m.Subcommands()
			if err != nil {
				errs = append(errs, err)
			}

			fcmds, ferrs := expand(subcmds, depth-1, includeExpandedModules)
			cmds = append(cmds, fcmds...)
			errs = append(errs, ferrs...)

			return cmds, errs
		} else {
			return []Command{cmd}, []error{}
		}
	})
}

func parallelMap[T any, R any](inputs []T, fn func(T) ([]R, []error)) ([]R, []error) {
	var wg sync.WaitGroup

	results := make([]struct {
		outs []R
		errs []error
	}, len(inputs))

	for i, input := range inputs {
		wg.Add(1)
		go func(idx int, t T) {
			defer wg.Done()
			fouts, ferrs := fn(t)
			results[idx].errs = append(results[idx].errs, ferrs...)
			results[idx].outs = append(results[idx].outs, fouts...)
		}(i, input)
	}

	wg.Wait()

	var outs []R
	var errs []error
	for _, result := range results {
		outs = append(outs, result.outs...)
		errs = append(errs, result.errs...)
	}

	return outs, errs
}
