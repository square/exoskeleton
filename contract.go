package exoskeleton

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/square/exit"
)

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
						Message: fmt.Sprintf("error reading %s: %s", cmd.path, err),
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
	summary, err := getMessageFromCommand(cmd, "summary")

	if err != nil {
		return "",
			exit.Wrap(
				CommandSummaryError{
					CommandError{
						Message: fmt.Sprintf("error getting summary from %s: %s", cmd.path, err),
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
	help, err := getMessageFromCommand(cmd, "help")

	if err != nil {
		return "",
			exit.Wrap(
				CommandHelpError{
					CommandError{
						Message: fmt.Sprintf("error getting help from %s: %s", cmd.path, err),
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
	cmd.Stderr = nil
	output, err := cmd.Output()
	if err != nil {
		return &commandDescriptor{},
			exit.Wrap(
				CommandDescribeError{
					CommandError{
						Message: fmt.Sprintf("error executing `%s --describe-commands`: %s", m.path, err),
						Command: m,
						Cause:   err,
					},
				},
				exit.InternalError,
			)
	}

	var descriptor *commandDescriptor
	if err := json.Unmarshal(output, &descriptor); err != nil {
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

// detectType reads the first two bytes from a file.
// If they are `#!`, we can assume that the file is a shell script.
//
// Armed with this assumption, we can extract the command's documentation
// using the magic comments approach.
//
// If the command is not a shell script, we will have to execute it
// to request its documentation.
func detectType(f *os.File) (fileType, error) {
	buffer := make([]byte, 2)
	_, err := f.Read(buffer)

	if err != nil {
		return unknown, fmt.Errorf("detectType: %w", err)
	} else if string(buffer) == "#!" {
		return script, nil
	} else {
		return binary, nil
	}
}

type fileType int

const (
	script fileType = iota
	binary
	unknown
)

func getMessageFromCommand(cmd *executableCommand, message string) (string, error) {
	f, err := os.Open(cmd.path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	t, err := detectType(f)
	if err != nil {
		return "", err
	}

	switch t {
	case script:
		s, err := getMessageFromMagicComments(f, message)
		if s == "" {
			return getMessageFromExecution(cmd, message)
		} else {
			return s, err
		}
	case binary:
		return getMessageFromExecution(cmd, message)
	default:
		return "", fmt.Errorf("Invalid value for t: %v", t)
	}
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
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("failed to execute %s: %w", cmd.Path, err)
	}
	return strings.TrimRight(string(out), "\n"), err
}
