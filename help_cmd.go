package exoskeleton

import (
	"fmt"
	"regexp"

	"github.com/square/exit"
)

// HelpHelp is the help text for the built-in 'help' command.
const HelpHelp = `USAGE
   %[1]s help [<command>]

OPTIONS
   -a, --all   Expand submodules and display all subcommands

EXAMPLES
   Display the documentation for the 'foobar' command
   $ %[1]s help foobar

   Expand modules and display a menu of all commands available
   $ %[1]s --all`

// HelpExec implements the 'help' command.
func HelpExec(e *Entrypoint, args, _ []string) error {
	if cmd, rest, err := e.Identify(args); err != nil {
		return err
	} else if IsNull(cmd) {
		return exit.ErrUnknownSubcommand
	} else if help, err := e.helpFor(cmd, rest); err != nil {
		return err
	} else {
		printHelp(help)
		return nil
	}
}

func (e *Entrypoint) helpFor(cmd Command, args []string) (string, error) {
	if m, ok := cmd.(Module); ok {
		return e.buildModuleHelp(m, args)
	} else {
		return cmd.Help()
	}
}

func (e *Entrypoint) printModuleHelp(m Module, args []string) error {
	help, err := e.buildModuleHelp(m, args)
	printHelp(help)
	return err
}

func (e *Entrypoint) buildModuleHelp(m Module, args []string) (string, error) {
	// If `Subcommands()` will return an error, return early
	if _, err := m.Subcommands(); err != nil {
		return "", err
	}

	cache := &summaryCache{Path: e.cachePath, onError: e.onError}
	opts := &MenuOptions{
		HeadingFor: e.menuHeadingFor,
		SummaryFor: cache.Read,
		Template:   e.menuTemplate,
	}

	for _, arg := range args {
		if arg == "--" {
			break
		} else if arg == "--all" || arg == "-a" {
			opts.Depth = -1
		}
	}

	menu, errs := MenuFor(m, opts)
	for _, err := range errs {
		e.onError(err)
	}
	return menu, nil
}

func printHelp(help string) {
	fmt.Println(formatHelp(help))
	fmt.Println()
}

func formatHelp(help string) string {
	re := regexp.MustCompile(`(?m)^([A-Z ]+)$`)
	return string(re.ReplaceAll([]byte(help), []byte("\033[1m$1\033[0m")))
}
