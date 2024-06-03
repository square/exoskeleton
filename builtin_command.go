package exoskeleton

import (
	"github.com/square/exoskeleton/pkg/shellcomp"
)

type builtinCommand struct {
	parent     Module
	definition *EmbeddedCommand
}

func (c *builtinCommand) Parent() Module  { return c.parent }
func (c *builtinCommand) Path() string    { return c.parent.Path() }
func (c *builtinCommand) Name() string    { return c.definition.Name }
func (c *builtinCommand) Summary() string { return c.definition.Summary }
func (c *builtinCommand) Help() string    { return c.definition.Help }

func (c *builtinCommand) Exec(e *Entrypoint, args, env []string) error {
	return c.definition.Exec(e, args, env)
}

func (c *builtinCommand) Complete(e *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error) {
	if c.definition.Complete != nil {
		return c.definition.Complete(e, args, env)
	} else {
		return []string{}, shellcomp.DirectiveNoFileComp, nil
	}
}
