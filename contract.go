package exoskeleton

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

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
			return getMessageFromExecution(cmd.Command("--" + message))
		} else {
			return s, err
		}
	case binary:
		return getMessageFromExecution(cmd.Command("--" + message))
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

func getMessageFromExecution(cmd *exec.Cmd) (string, error) {
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("failed to execute %s: %w", cmd.Path, err)
	}
	return strings.TrimRight(string(out), "\n"), err
}

func getMessageFromDir(path string, flag string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return getMessageFromMagicComments(f, flag)
}
