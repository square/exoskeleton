package exoskeleton

import (
	"strings"

	"github.com/square/exoskeleton/pkg/shellcomp"
)

// CompleteCommands is a CompleteFunc that provides completions for command names.
// It is used by commands like 'help' and 'which' which expect their arguments
// to be command names.
func CompleteCommands(e *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error) {
	return e.completionsFor(args, env, false)
}

// CompleteFiles is a CompleteFunc that provides completions for files and paths.
func CompleteFiles(e *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error) {
	return nil, shellcomp.DirectiveDefault, nil
}

func (e *Entrypoint) completionsFor(args, env []string, completeArgs bool) ([]string, shellcomp.Directive, error) {
	if len(args) == 0 {
		return nil, shellcomp.DirectiveNoFileComp, nil
	}

	// The last argument, which is not completely typed by the user,
	// should not be part of the list of arguments
	toComplete := args[len(args)-1]
	trimmedArgs := args[:len(args)-1]

	// Find the real command for which completion must be performed
	finalCmd, finalCmdArgs, err := e.Identify(trimmedArgs)
	if err != nil {
		return nil, shellcomp.DirectiveError, err
	}

	if _, isModule := finalCmd.(Module); !isModule && !completeArgs {
		return nil, shellcomp.DirectiveNoFileComp, nil
	}

	return finalCmd.Complete(e, append(finalCmdArgs, toComplete), env)
}

func completionsForModule(m Module, args []string) ([]string, shellcomp.Directive, error) {
	cmds, err := m.Subcommands()
	if err != nil {
		return nil, shellcomp.DirectiveError, err
	}
	return cmds.completionsFor(args)
}

func (c Commands) completionsFor(args []string) ([]string, shellcomp.Directive, error) {
	var completions []string

	if len(args) > 0 {
		toComplete := args[0]
		var name string

		seen := make(map[string]bool)

		for _, subcmd := range c {
			name = subcmd.Name()

			if seen[name] {
				continue
			}
			seen[name] = true

			if strings.HasPrefix(name, toComplete) {
				completions = append(completions, name)
			}
		}
	}

	return completions, shellcomp.DirectiveNoFileComp, nil
}
