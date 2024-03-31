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
	} else if p, ok := cmd.(helpWithArgsProvider); ok {
		return p.helpWithArgs(args)
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
	cmds, err := m.Subcommands()
	if err != nil {
		return "", err
	}

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
		cmds = e.expandModules(cmds)
	}

	return e.MenuRelativeTo(cmds, m).String(), nil
}

func (e *Entrypoint) expandModules(cmds Commands) Commands {
	all := Commands{}
	for _, cmd := range cmds {
		if m, ok := cmd.(Module); ok {
			subcmds, err := m.Subcommands()
			if err != nil {
				e.onError(err)
			}
			all = append(all, e.expandModules(subcmds)...)
		} else {
			all = append(all, cmd)
		}
	}
	return all
}

func printHelp(help string) {
	fmt.Println(formatHelp(help))
	fmt.Println()
}

func formatHelp(help string) string {
	re := regexp.MustCompile(`(?m)^([A-Z ]+)$`)
	return string(re.ReplaceAll([]byte(help), []byte("\033[1m$1\033[0m")))
}
