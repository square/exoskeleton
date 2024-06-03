package exoskeleton

import (
	"errors"
	"os"
	"os/exec"

	"github.com/square/exit"
	"github.com/square/exoskeleton/pkg/shellcomp"
)

// executableCommand implements the Command interface for a file that can be executed.
type executableCommand struct {
	entrypoint   *Entrypoint
	parent       Module
	path         string
	name         string
	args         []string
	summary      string
	discoveredIn string
}

func (cmd *executableCommand) Parent() Module       { return cmd.parent }
func (cmd *executableCommand) Path() string         { return cmd.path }
func (cmd *executableCommand) Name() string         { return cmd.name }
func (cmd *executableCommand) DiscoveredIn() string { return cmd.discoveredIn }

// Command returns an exec.Cmd that will run the executable with the given arguments.
func (cmd *executableCommand) Command(args ...string) *exec.Cmd {
	return exec.Command(cmd.path, append(cmd.args, args...)...)
}

// Exec invokes the executable with the given arguments and environment.
func (cmd *executableCommand) Exec(_ *Entrypoint, args, env []string) error {
	command := cmd.Command(args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = env

	err := command.Run()
	if err == nil {
		return nil
	}

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return exit.Wrap(nil, exitError.ExitCode())
	}

	return err
}

// Complete invokes the executable with `--complete` as its first argument
// and parses its output according to Cobra's ShellComp API.
func (cmd *executableCommand) Complete(_ *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error) {
	command := cmd.Command(append([]string{"--complete", "--"}, args...)...)
	command.Stdin = nil
	command.Stderr = nil
	command.Env = env

	if out, err := command.Output(); err != nil {
		return []string{}, shellcomp.DirectiveNoFileComp, err
	} else {
		return shellcomp.Unmarshal(out)
	}
}

// Summary returns the (short!) description of the command to be displayed
// in menus.
//
// When Command is a shell script, it reads the script's source to extract a
// single-line comment like '# SUMMARY: summary goes here'.
//
// When Command is a binary, it executes the command with the flag '--summary'.
// The executable is expected to write the summary to standard output and exit
// successfully.
func (cmd *executableCommand) Summary() string {
	if cmd.summary != "" {
		return cmd.summary
	}

	return cmd.entrypoint.readSummaryFromExecutable(cmd)
}

// Help returns the help text for the command.
//
// When Command is a shell script, it reads the script's source to extract a
// multi-line magic comment that starts with '# HELP:'.
//
// When Command is a binary, it executes the command with the flag '--help'.
// The executable is expected to write the help text to standard output and exit
// successfully.
func (cmd *executableCommand) Help() string {
	return cmd.helpWithArgs(nil)
}

// helpWithArgs returns the help text for the command, passing
// any arguments through to the executable that responds to --help.
func (cmd *executableCommand) helpWithArgs(args []string) string {
	return cmd.entrypoint.readHelpFromExecutable(cmd, args)
}

// helpWithArgsProvider can be implemented by commands to
// provide HelpWithArgs, to accept command line arguments
// when printing help.
//
// It is not exported yet because its API isn't stable.
type helpWithArgsProvider interface {
	helpWithArgs([]string) string
}
