package exoskeleton

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/square/exit"
	"github.com/square/exoskeleton/pkg/shellcomp"
)

// Contract represents an agreement between Exoskeleton and a command about
// how the command is discovered.
//
// Contracts are tried in order during discovery. The first contract that
// doesn't return ErrNotApplicable builds the command.
type Contract interface {
	// BuildCommand constructs a Command using this contract's rules.
	// Returns ErrNotApplicable if this contract doesn't handle the file/directory.
	// Returns nil, nil if the file/directory should be ignored (e.g., non-executable file).
	BuildCommand(path string, info fs.DirEntry, parent Module, d DiscoveryContext) (Command, error)
}

// ErrNotApplicable indicates that a contract does not apply to a given file/directory.
// Discovery will try the next contract in the list.
var ErrNotApplicable = errors.New("contract does not apply")

// CommandError records an error that occurred with a command's implementation of its interface
type CommandError struct {
	Message string
	Cause   error
	Command Command
}

func (e CommandError) Error() string { return e.Message }
func (e CommandError) Unwrap() error { return e.Cause }

// CommandSummaryError indicates that a command did not properly implement the
// interface for providing a summary
type CommandSummaryError struct{ CommandError }

// CommandHelpError indicates that a command did not properly implement the
// interface for providing help
type CommandHelpError struct{ CommandError }

// CommandDescribeError indicates that an executable module did not properly
// respond to `--describe-commands`
type CommandDescribeError struct{ CommandError }

func readSummaryFromModulefile(cmd *directoryModule) (string, error) {
	var summary string

	f, err := os.Open(cmd.path)
	if err == nil {
		defer f.Close()
		summary, err = getMessageFromMagicComments(f, "summary")
	}

	if err != nil {
		return "",
			exit.Wrap(
				CommandSummaryError{
					CommandError{
						Message: fmt.Sprintf("summary('%s'): %s", Usage(cmd), err),
						Command: cmd,
						Cause:   err,
					},
				},
				exit.InternalError,
			)
	}

	return summary, nil
}

func readSummaryFromExecutable(cmd *executableCommand) (string, error) {
	summary, err := getMessageFromExecution(cmd, "summary")

	if err != nil {
		return "",
			exit.Wrap(
				CommandSummaryError{
					CommandError{
						Message: fmt.Sprintf("summary('%s'): %s", Usage(cmd), err),
						Command: cmd,
						Cause:   err,
					},
				},
				exit.InternalError,
			)
	}

	return summary, nil
}

func readHelpFromExecutable(cmd *executableCommand) (string, error) {
	help, err := getMessageFromExecution(cmd, "help")

	if err != nil {
		return "",
			exit.Wrap(
				CommandHelpError{
					CommandError{
						Message: fmt.Sprintf("help('%s'): %s", Usage(cmd), err),
						Command: cmd,
						Cause:   err,
					},
				},
				exit.InternalError,
			)
	}

	return help, nil
}

func describeCommands(m *executableModule) (*commandDescriptor, error) {
	cmd := m.Command("--describe-commands")
	out, err := m.output(cmd)
	if err != nil {
		err = fmt.Errorf("exec '%s': %w", strings.Join(cmd.Args, " "), err)

		return &commandDescriptor{},
			exit.Wrap(
				CommandDescribeError{
					CommandError{
						Message: err.Error(),
						Command: m,
						Cause:   err,
					},
				},
				exit.InternalError,
			)
	}

	var descriptor *commandDescriptor
	if err := json.Unmarshal(out, &descriptor); err != nil {
		return &commandDescriptor{},
			exit.Wrap(
				CommandDescribeError{
					CommandError{
						Message: fmt.Sprintf("error parsing output from `%s --describe-commands`: %s", m.path, err),
						Command: m,
						Cause:   err,
					},
				},
				exit.InternalError,
			)
	}

	return descriptor, nil
}

type commandDescriptor struct {
	Name     string               `json:"name"`
	Summary  *string              `json:"summary,omitempty"`
	Commands []*commandDescriptor `json:"commands,omitempty"`
}

func readSummaryFromShellScript(cmd *shellScriptCommand) (string, error) {
	f, err := os.Open(cmd.path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	summary, err := getMessageFromMagicComments(f, "summary")
	if err != nil {
		return "",
			exit.Wrap(
				CommandSummaryError{
					CommandError{
						Message: fmt.Sprintf("summary('%s'): %s", Usage(cmd), err),
						Command: cmd,
						Cause:   err,
					},
				},
				exit.InternalError,
			)
	}

	// Fall back to executing with --summary if magic comments are empty
	if summary == "" {
		return getMessageFromExecution(&cmd.executableCommand, "summary")
	}

	return summary, nil
}

func readHelpFromShellScript(cmd *shellScriptCommand) (string, error) {
	f, err := os.Open(cmd.path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	help, err := getMessageFromMagicComments(f, "help")
	if err != nil {
		return "",
			exit.Wrap(
				CommandHelpError{
					CommandError{
						Message: fmt.Sprintf("help('%s'): %s", Usage(cmd), err),
						Command: cmd,
						Cause:   err,
					},
				},
				exit.InternalError,
			)
	}

	// Fall back to executing with --help if magic comments are empty
	if help == "" {
		return getMessageFromExecution(&cmd.executableCommand, "help")
	}

	return help, nil
}

func getMessageFromMagicComments(f *os.File, message string) (string, error) {
	reader := bufio.NewReader(f)
	switch message {
	case "summary":
		return getSummaryFromMagicComments(reader)
	case "help":
		return getHelpFromMagicComments(reader)
	default:
		panic("Unhandled message: " + message)
	}
}

func getSummaryFromMagicComments(reader *bufio.Reader) (string, error) {
	var line string
	var err error

	for {
		line, err = reader.ReadString('\n')
		if strings.HasPrefix(line, "# SUMMARY:") {
			return strings.TrimRight(strings.TrimPrefix(line[10:], " "), "\n"), nil
		}
		if err == io.EOF {
			return "", nil
		}
		if err != nil {
			return "", err
		}
	}
}

func getHelpFromMagicComments(reader *bufio.Reader) (string, error) {
	var line string
	var err error
	var help string
	var inHelpText bool

	for {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			return strings.TrimRight(help, "\n"), nil
		}

		if err != nil {
			return "", err
		}

		if strings.HasPrefix(line, "# USAGE:") {
			help += "USAGE\n   " + strings.TrimRight(strings.TrimPrefix(line[8:], " "), "\n") + "\n\n"
		}

		if inHelpText {
			if strings.HasPrefix(line, "#") {
				if len(line) > 2 {
					help += line[2:]
				} else {
					help += "\n"
				}
			} else {
				inHelpText = false
			}
		}

		if strings.HasPrefix(line, "# HELP:") {
			help += strings.TrimPrefix(line[7:], " ")
			inHelpText = true
		}
	}
}

func getMessageFromExecution(c *executableCommand, message string) (string, error) {
	cmd := c.Command("--" + message)
	out, err := c.output(cmd)
	if err != nil {
		err = fmt.Errorf("exec '%s': %w", strings.Join(cmd.Args, " "), err)
	}
	return strings.TrimRight(string(out), "\n"), err
}

func getCompletionsFromExecutable(c *executableCommand, args, env []string) ([]string, shellcomp.Directive, error) {
	cmd := c.Command(append([]string{"--complete", "--"}, args...)...)
	cmd.Env = env

	out, err := c.output(cmd)
	if err != nil {
		return []string{}, shellcomp.DirectiveNoFileComp, err
	}

	return shellcomp.Unmarshal(out)
}
