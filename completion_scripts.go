package exoskeleton

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/spf13/cobra"
)

// GenerateCompletionScript generates a completion script for the given shell ("bash" or "zsh")
// and writes it to the given writer.
func (e *Entrypoint) GenerateCompletionScript(shell string, w io.Writer) error {
	return GenerateCompletionScript(e.name, shell, w)
}

// GenerateCompletionScript generates a completion script for the given shell ("bash" or "zsh")
// and writes it to the given writer.
func GenerateCompletionScript(name, shell string, w io.Writer) (err error) {
	c := &cobra.Command{Use: name}
	b := new(bytes.Buffer)

	if shell == "bash" {
		err = c.GenBashCompletionV2(b, true)
	} else if shell == "zsh" {
		err = c.GenZshCompletion(b)
	} else {
		err = fmt.Errorf("unsupported shell: %s", shell)
	}

	if err != nil {
		return err
	}

	// Cobra CLIs generate completions in response to a command
	// named "__complete". Exoskeleton names the command "complete" instead
	re := regexp.MustCompile(`\b__complete\b`)
	_, err = w.Write([]byte(re.ReplaceAllString(b.String(), "complete")))
	return err
}
