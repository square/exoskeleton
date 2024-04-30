package exoskeleton

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/square/exit"
	"github.com/square/exoskeleton/pkg/shellcomp"
)

// executableCommand implements the Command interface for a file that can be executed.
type executableCommand struct {
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
	c := exec.Command(cmd.path, append(cmd.args, args...)...)
	fmt.Fprintf(os.Stderr, "\033[90m%s\033[0m\n", strings.Join(c.Args, " "))
	return c
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
func (cmd *executableCommand) Summary() (string, error) {
	if cmd.summary != "" {
		return cmd.summary, nil
	}

	return getMessageFromCommand(cmd, "summary")
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
	return getMessageFromCommand(cmd, "help")
}

// Reads the first two bytes from a file.
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
			return cmd.getMessageFromExecution(message)
		} else {
			return s, err
		}
	case binary:
		return cmd.getMessageFromExecution(message)
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
			return help, nil
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

func (cmd *executableCommand) getMessageFromExecution(flag string) (string, error) {
	command := cmd.Command("--" + flag)
	command.Stderr = nil
	out, err := command.Output()
	if err != nil {
		err = fmt.Errorf("failed to execute %s: %w", cmd.path, err)
	}
	return strings.TrimRight(string(out), "\n"), err
}
