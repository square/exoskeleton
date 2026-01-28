package exoskeleton

import (
	"fmt"
	"os"

	"github.com/square/exit"
	"github.com/square/exoskeleton/v2/pkg/shellcomp"
)

// CompleteHelp is the help text for the built-in 'complete' command.
const CompleteHelp = `USAGE
   %[1]s complete <prefix>

   Provides a list of completions for <prefix> followed by a completion directive.
   Used by Bash and Zsh completion scripts.

EXAMPLES
   List 'help' and any other commands that start with 'hel'
   $ %[1]s complete hel

COMPLETION DIRECTIVES
   The directive is a bitmap that combines one or more of the following behaviors:

    0  Default                    Let the shell perform its default behavior after providing a completion
    1  Error                      An error occurred and completions should be ignored
    2  No Space                   The shell should not add a space after providing a completion
    4  No File Completions        The shell should not provide file completions if no completions are listed
    8  Filter Files by Extension  The shell should use the provided completions as file extension filters
   16  Directories                The shell should provide file completions but only suggest directories`

// CompleteExec implements the 'complete' command.
func CompleteExec(e *Entrypoint, args, env []string) error {
	completions, directive, err := e.completionsFor(args, env, true)

	if err != nil {
		completionError(err.Error())
		// Keep going for multiple reasons:
		// 1) There could be some valid completions even though there was an error
		// 2) Even without completions, we need to print the directive
	}

	os.Stdout.Write(shellcomp.Marshal(completions, directive, false))

	// Print some helpful info to stderr for the user to understand.
	// Output from stderr must be ignored by the completion script.
	fmt.Fprintf(os.Stderr, "Completion ended with directive: %s\n", directive)

	return nil
}

// completionError prints the specified completion message to stderr.
func completionError(s string) {
	s = fmt.Sprintf("[Error] %s\n", s)
	completionDebug(s)

	// Note that completion printouts should never be on stdout as they would
	// be wrongly interpreted as actual completion choices by the completion script.
	fmt.Fprint(os.Stderr, s)
}

// completionDebug prints the specified string to the same file as where the
// completion script prints its logs.
func completionDebug(s string) {
	// Such logs are only printed when the user has set the environment
	// variable BASH_COMP_DEBUG_FILE to the path of some file to be used.
	if path := os.Getenv("BASH_COMP_DEBUG_FILE"); path != "" {
		if f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			defer f.Close()

			// Write msg to f. If there's an error, exit.
			if _, err := f.WriteString(s); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(exit.NotOK)
			}
		}
	}
}
