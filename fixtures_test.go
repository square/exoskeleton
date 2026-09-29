package exoskeleton

import (
	"path/filepath"
	"runtime"

	"github.com/square/exoskeleton/v2/pkg/shellcomp"
)

var fixtures string

func init() {
	_, testfile, _, _ := runtime.Caller(0)
	fixtures = filepath.Join(testfile, "..", "fixtures")
}

// stubParent is a Command that reports that it may have subcommands
// and records whether Subcommands() was ever called.
type stubParent struct {
	name              string
	parent            Command
	subcommandsCalled bool
}

func (c *stubParent) Path() string             { return "" }
func (c *stubParent) Name() string             { return c.name }
func (c *stubParent) Parent() Command          { return c.parent }
func (c *stubParent) Aliases() []string        { return nil }
func (c *stubParent) Summary() (string, error) { return "A stub", nil }
func (c *stubParent) Help() (string, error)    { panic("Unused") }
func (c *stubParent) HasSubcommands() bool     { return true }

func (c *stubParent) Exec(*Entrypoint, []string, []string) error { panic("Unused") }

func (c *stubParent) Complete(*Entrypoint, []string, []string) ([]string, shellcomp.Directive, error) {
	panic("Unused")
}

func (c *stubParent) DefaultSubcommand() Command { return nil }

func (c *stubParent) Subcommands() (Commands, error) {
	c.subcommandsCalled = true
	return Commands{}, nil
}
