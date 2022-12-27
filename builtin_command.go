package exoskeleton

import (
	"github.com/square/exoskeleton/pkg/shellcomp"
)

type builtinCommand struct {
	parent   Module
	name     string
	summary  string
	help     string
	exec     ExecFunc
	complete CompleteFunc
}

// ExecFunc is called when an built-in command is run.
type ExecFunc func(e *Entrypoint, args, env []string) error

// CompleteFunc is called when an built-in command is asked to supply shell completions.
type CompleteFunc func(e *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error)

func (c *builtinCommand) Parent() Module           { return c.parent }
func (c *builtinCommand) Path() string             { return c.parent.Path() }
func (c *builtinCommand) Name() string             { return c.name }
func (c *builtinCommand) Summary() (string, error) { return c.summary, nil }
func (c *builtinCommand) Help() (string, error)    { return c.help, nil }

func (c *builtinCommand) Exec(e *Entrypoint, args, env []string) error {
	return c.exec(e, args, env)
}

func (c *builtinCommand) Complete(e *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error) {
	if c.complete != nil {
		return c.complete(e, args, env)
	} else {
		return []string{}, shellcomp.DirectiveNoFileComp, nil
	}
}
