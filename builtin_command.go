package exoskeleton

import (
	"github.com/square/exoskeleton/pkg/shellcomp"
)

type builtinCommand struct {
	parent      Command
	definition  *EmbeddedCommand
	subcommands Commands // Empty for leaf commands
}

func (c *builtinCommand) Parent() Command          { return c.parent }
func (c *builtinCommand) Path() string             { return c.parent.Path() }
func (c *builtinCommand) Name() string             { return c.definition.Name }
func (c *builtinCommand) Summary() (string, error) { return c.definition.Summary, nil }
func (c *builtinCommand) Help() (string, error)    { return c.definition.Help, nil }

func (c *builtinCommand) Exec(e *Entrypoint, args, env []string) error {
	if len(c.subcommands) > 0 {
		return e.printModuleHelp(c, args)
	}
	return c.definition.Exec(e, args, env)
}

func (c *builtinCommand) Complete(e *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error) {
	if len(c.subcommands) > 0 {
		return completionsForSubcommands(c, args)
	}
	if c.definition.Complete != nil {
		return c.definition.Complete(e, args, env)
	}
	return []string{}, shellcomp.DirectiveNoFileComp, nil
}

func (c *builtinCommand) Subcommands() (Commands, error) {
	return c.subcommands, nil
}
