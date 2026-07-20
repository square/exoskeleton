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

// WhichExec implements the 'which' command.
func WhichExec(e *Entrypoint, args, _ []string) error {
	identifyArgs, willResolveSymlinks := splitWhichArgs(args)

	cmd, _, err := e.Identify(identifyArgs)
	if err != nil {
		return err
	} else if IsNull(cmd) {
		return exit.ErrUnknownSubcommand
	}

	path := cmd.Path()
	if willResolveSymlinks {
		resolvedPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unable to follow symlinks in %s\n", path)
			return err
		}
		path = resolvedPath
	}

	fmt.Println(path)
	return nil
}

// splitWhichArgs separates which's own --follow-symlinks/-s flags from the
// arguments that identify the command. The flags are consumed here rather than
// forwarded to Identify: left in, a flag becomes a trailing argument and causes
// Identify to resolve the command's default subcommand instead of the command
// itself. Everything after a `--` terminator is left untouched for the command.
func splitWhichArgs(args []string) (identifyArgs []string, followSymlinks bool) {
	identifyArgs = make([]string, 0, len(args))

	for i, arg := range args {
		if arg == "--" {
			identifyArgs = append(identifyArgs, args[i:]...)
			break
		} else if arg == "--follow-symlinks" || arg == "-s" {
			followSymlinks = true
		} else {
			identifyArgs = append(identifyArgs, arg)
		}
	}

	return identifyArgs, followSymlinks
}
