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

// helpExec implements the 'help' command.
func helpExec(e *Entrypoint, args, _ []string) error {
	if cmd, rest := e.Identify(args); IsNull(cmd) {
		return exit.ErrUnknownSubcommand
	} else {
		help := e.helpFor(cmd, rest)
		printHelp(help)
		return nil
	}
}

func (e *Entrypoint) helpFor(cmd Command, args []string) string {
	if m, ok := cmd.(Module); ok {
		return e.buildModuleHelp(m, args)
	} else if p, ok := cmd.(helpWithArgsProvider); ok {
		return p.helpWithArgs(args)
	} else {
		return cmd.Help()
	}
}

func (e *Entrypoint) printModuleHelp(m Module, args []string) error {
	printHelp(e.buildModuleHelp(m, args))
	return nil
}

func (e *Entrypoint) buildModuleHelp(m Module, args []string) string {
	cmds := m.Subcommands()

	var filteredArgs []string
	var willExpandMenu bool

	for i, arg := range args {
		if arg == "--" {
			filteredArgs = append(filteredArgs, args[i:]...)
			break
		} else if arg == "--all" || arg == "-a" {
			willExpandMenu = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	if willExpandMenu {
		cmds = withoutModules(cmds.Flatten())
	}

	return e.buildMenu(cmds, m).String()
}

func withoutModules(cmds Commands) (all []Command) {
	for _, c := range cmds {
		if _, ok := c.(Module); !ok {
			all = append(all, c)
		}
	}
	return
}

func printHelp(help string) {
	fmt.Println(formatHelp(help))
	fmt.Println()
}

func formatHelp(help string) string {
	re := regexp.MustCompile(`(?m)^([A-Z ]+)$`)
	return string(re.ReplaceAll([]byte(help), []byte("\033[1m$1\033[0m")))
}
