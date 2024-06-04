package exoskeleton

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/square/exit"
)

// WhichHelp is the help text for the built-in 'which' command.
const WhichHelp = `USAGE
   %[1]s which [<command>]

   Displays the path where the given command exists.
   Displays the path to %[1]s for built-in commands like which and help.

OPTIONS
   -s, --follow-symlinks   Follow symlinks before displaying the path

EXAMPLES
   %[1]s which            # Display the path to %[1]s
   %[1]s which help       # Display the path to %[1]s
   %[1]s which foobar     # Display the path to the foobar command`

// whichExec implements the 'which' command.
func whichExec(e *Entrypoint, args, _ []string) error {
	if cmd, _, err := e.Identify(args); err != nil {
		return err
	} else if IsNull(cmd) {
		return exit.ErrUnknownSubcommand
	} else {
		willResolveSymlinks := false

		for _, arg := range args {
			if arg == "--" {
				break
			} else if arg == "--follow-symlinks" || arg == "-s" {
				willResolveSymlinks = true
			}
		}

		path := cmd.Path()
		if willResolveSymlinks {
			resolvedPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Unable to follow symlinks in %s\n", path)
				return err
			} else {
				path = resolvedPath
			}
		}

		fmt.Println(path)
		return nil
	}
}
