// The constants in this file are taken from https://github.com/spf13/cobra/blob/v1.5.0/completions.go
//
// Cobra defines a completion API so that generic completion logic can be written once for Bash and Zsh
// that invokes a CLI with `--complete --` and passes the inputs to be completed.
//
// The CLI is expected to write suggestions to Standard Output as well as a Shell Completion Directive,
// which is an integer that stacks one or more bit flags.
package shellcomp

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Directive is a bit map representing the different behaviors the shell
// can be instructed to have once completions have been provided.
type Directive int

const (
	// DirectiveError indicates an error occurred and completions should be ignored.
	DirectiveError Directive = 1 << iota

	// DirectiveNoSpace indicates that the shell should not add a space
	// after the completion even if there is a single completion provided.
	// Used when a flag ends with '='.
	DirectiveNoSpace

	// DirectiveNoFileComp indicates that the shell should not provide
	// file completion even when no completion is provided.
	DirectiveNoFileComp

	// DirectiveFilterFileExt indicates that the provided completions
	// should be used as file extension filters.
	// For flags, using Command.MarkFlagFilename() and Command.MarkPersistentFlagFilename()
	// is a shortcut to using this directive explicitly. The BashCompFilenameExt
	// annotation can also be used to obtain the same behavior for flags.
	DirectiveFilterFileExt

	// DirectiveFilterDirs indicates that only directory names should
	// be provided in file completion. To request directory names within another
	// directory, the returned completions should specify the directory within
	// which to search. The BashCompSubdirsInDir annotation can be used to
	// obtain the same behavior but only for flags.
	DirectiveFilterDirs

	// ===========================================================================

	// All directives using iota should be above this one.
	// For internal use.
	directiveMaxValue

	// DirectiveDefault indicates to let the shell perform its default
	// behavior after completions have been provided.
	// This one must be last to avoid messing up the iota count.
	DirectiveDefault Directive = 0
)

// Returns a string listing the different directive enabled in the specified parameter
func (d Directive) String() string {
	var directives []string
	if d&DirectiveError != 0 {
		directives = append(directives, "ShellCompDirectiveError")
	}
	if d&DirectiveNoSpace != 0 {
		directives = append(directives, "ShellCompDirectiveNoSpace")
	}
	if d&DirectiveNoFileComp != 0 {
		directives = append(directives, "ShellCompDirectiveNoFileComp")
	}
	if d&DirectiveFilterFileExt != 0 {
		directives = append(directives, "ShellCompDirectiveFilterFileExt")
	}
	if d&DirectiveFilterDirs != 0 {
		directives = append(directives, "ShellCompDirectiveFilterDirs")
	}
	if len(directives) == 0 {
		directives = append(directives, "ShellCompDirectiveDefault")
	}

	if d >= directiveMaxValue {
		return fmt.Sprintf("ERROR: unexpected ShellCompDirective value: %d", d)
	}
	return strings.Join(directives, ", ")
}

func Marshal(completions []string, directive Directive, noDescriptions bool) (result []byte) {
	var b bytes.Buffer

	for _, comp := range completions {
		if noDescriptions {
			// Remove any description that may be included following a tab character.
			comp = strings.Split(comp, "\t")[0]
		}

		// Make sure we only write the first line to the output.
		// This is needed if a description contains a linebreak.
		// Otherwise the shell scripts will interpret the other lines as new flags
		// and could therefore provide a wrong completion.
		comp = strings.Split(comp, "\n")[0]

		// Finally trim the completion. This is especially important to get rid
		// of a trailing tab when there are no description following it.
		// For example, a sub-command without a description should not be completed
		// with a tab at the end (or else zsh will show a -- following it
		// although there is no description).
		comp = strings.TrimSpace(comp)

		// Print each possible completion to stdout for the completion script to consume.
		fmt.Fprintln(&b, comp)
	}

	// As the last printout, print the completion directive for the completion script to parse.
	// The directive integer must be that last character following a single colon (:).
	// The completion script expects :<directive>
	fmt.Fprintf(&b, ":%d\n", directive)

	return b.Bytes()
}

func Unmarshal(bytes []byte) ([]string, Directive, error) {
	lines := strings.Split(strings.TrimSuffix(string(bytes), "\n"), "\n")

	directiveString := strings.TrimPrefix(lines[len(lines)-1], ":") // ":4"
	completions := lines[:len(lines)-1]

	if directive, err := strconv.Atoi(directiveString); err != nil {
		return completions, DirectiveNoFileComp, err
	} else {
		return completions, Directive(directive), nil
	}
}
