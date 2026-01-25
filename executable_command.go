package exoskeleton

import (
	"bytes"
	"errors"
	"os"
	"os/exec"

	"github.com/square/exoskeleton/pkg/shellcomp"
)

// executableCommand implements the Command interface for a file that can be executed.
type executableCommand struct {
	parent       Command
	path         string
	name         string
	args         []string
	summary      *string
	discoveredIn string
	executor     ExecutorFunc
	cmds         Commands
	discoverer   DiscoveryContext
	cache        Cache
}

func (cmd *executableCommand) Parent() Command      { return cmd.parent }
func (cmd *executableCommand) Path() string         { return cmd.path }
func (cmd *executableCommand) Name() string         { return cmd.name }
func (cmd *executableCommand) DiscoveredIn() string { return cmd.discoveredIn }

// Command returns an exec.Cmd that will run the executable with the given arguments.
func (cmd *executableCommand) Command(args ...string) *exec.Cmd {
	return exec.Command(cmd.path, append(cmd.args, args...)...)
}

// Exec invokes the executable with the given arguments and environment.
// If this is a module (has discoverer), it prints the module help instead.
func (cmd *executableCommand) Exec(e *Entrypoint, args, env []string) error {
	if cmd.discoverer != nil {
		return e.printModuleHelp(cmd, args)
	}
	command := cmd.Command(args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = env
	return cmd.run(command)
}

// Complete invokes the executable with `--complete` as its first argument
// and parses its output according to Cobra's ShellComp API.
func (cmd *executableCommand) Complete(_ *Entrypoint, args, env []string) ([]string, shellcomp.Directive, error) {
	if cmd.discoverer != nil {
		return completionsForSubcommands(cmd, args)
	}
	return getCompletionsFromExecutable(cmd, args, env)
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
func (cmd *executableCommand) Summary() (string, error) {
	if cmd.discoverer != nil && cmd.cmds == nil {
		if err := cmd.discover(); err != nil {
			return "", err
		}
	}

	if cmd.summary != nil {
		return *cmd.summary, nil
	}

	return cmd.cache.Fetch(cmd, "summary", func() (string, error) {
		return readSummaryFromExecutable(cmd)
	})
}

// Help returns the help text for the command.
//
// When Command is a shell script, it reads the script's source to extract a
// multi-line magic comment that starts with '# HELP:'.
//
// When Command is a binary, it executes the command with the flag '--help'.
// The executable is expected to write the help text to standard output and exit
// successfully.
func (cmd *executableCommand) Help() (string, error) {
	return readHelpFromExecutable(cmd)
}

// Subcommands returns the list of subcommands for this command.
// Returns an empty slice for leaf commands.
func (cmd *executableCommand) Subcommands() (Commands, error) {
	if cmd.discoverer == nil {
		return Commands{}, nil // Leaf command
	}
	if cmd.cmds == nil {
		if err := cmd.discover(); err != nil {
			return nil, err
		}
	}
	return cmd.cmds, nil
}

func (cmd *executableCommand) run(c *exec.Cmd) error {
	return cmd.executor(c)
}

func (cmd *executableCommand) output(c *exec.Cmd) ([]byte, error) {
	// Expect to capture standard output and standard error
	if c.Stdout != nil {
		return nil, errors.New("exec: Stdout already set")
	}
	if c.Stderr != nil {
		return nil, errors.New("exec: Stderr already set")
	}
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := cmd.run(c)
	if ee, ok := err.(*exec.ExitError); ok {
		ee.Stderr = stderr.Bytes()
	}
	return stdout.Bytes(), err
}

// discover invokes an executable with `--describe-commands` and constructs a tree
// of modules and subcommands (all to be invoked through the given executable)
// from the JSON output.
func (cmd *executableCommand) discover() error {
	out, err := cmd.cache.Fetch(cmd, "describe-commands", func() (string, error) {
		return describeCommandsRaw(cmd)
	})

	if err != nil {
		return err
	}

	descriptor, err := parseDescribeCommands(cmd, out)
	if err != nil {
		return err
	}

	cmd.summary = descriptor.Summary
	cmd.cmds = toCommands(cmd, descriptor.Commands, nil, cmd.discoverer)
	return nil
}

func toCommands(parent *executableCommand, descriptors []*commandDescriptor, args []string, d DiscoveryContext) Commands {
	cmds := Commands{}
	for _, descriptor := range descriptors {
		c := &executableCommand{
			parent:       parent,
			discoveredIn: parent.discoveredIn,
			path:         parent.path,
			args:         append(args, descriptor.Name),
			name:         descriptor.Name,
			summary:      descriptor.Summary,
			executor:     parent.executor,
			cache:        parent.cache,
		}

		if len(descriptor.Commands) > 0 && d.MaxDepth() != 0 {
			c.discoverer = d.Next()
			c.cmds = toCommands(c, descriptor.Commands, append(args, c.name), d.Next())
		}
		cmds = append(cmds, c)
	}
	return cmds
}
