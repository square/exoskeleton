package exoskeleton

import "strings"

// Identify identifies the command being invoked.
//
// The function also returns any arguments that were not used to identify the command.
// For example, if the Go CLI were implemented with Exoskeleton, and we ran it
// with the rawArgs 'get -u github.com/square/exoskeleton/v2', Identify would return
// the Get command and the arguments {'-u', 'github.com/square/exoskeleton/v2'}.
//
// If no command is identified, Identify invokes CommandNotFound callbacks and
// returns NullCommand.
//
// Returns a CommandError if the command does not fulfill the contract
// for providing its subcommands.
func (e *Entrypoint) Identify(args []string) (Command, []string, error) {
	// Recognize `--complete` as an alias for the built-in `complete` command.
	if len(args) > 0 && args[0] == "--complete" {
		return e.Identify(append([]string{"complete"}, args[1:]...))
	}

	cmd, rest, err := identify(e, args)

	// Recognize `--help` and `-h` as aliases for the built-in `help` command
	// only when they immediately follow an identifiable command.
	if !IsNull(cmd) && len(rest) > 0 && (rest[0] == "--help" || rest[0] == "-h") {
		return e.Identify(append(append([]string{"help"}, argsRelativeTo(cmd, e)...), rest[1:]...))
	}

	if IsNull(cmd) {
		e.commandNotFound(cmd)
	} else if err == nil {
		e.afterIdentify(cmd, rest)
	}

	return cmd, rest, err
}

// identify uses args to identify a Command and returns the command and the rest
// of the commandline arguments or else {nil, args} if no Command is identified.
//
// Returns a CommandError if the command does not fulfill the contract
// for providing its subcommands.
func identify(cmd Command, args []string) (Command, []string, error) {
	if len(args) == 0 || isFlag(args[0]) {
		return cmd, args, nil
	}

	// Rewrite args like {"module:subcommand", "--flag"} to {"module", "subcommand", "--flag"}.
	// Do this just-in-time, non-destructively, while we're working on identifying a command.
	name, rest := args[0], args[1:]
	if strings.Contains(name, ":") {
		return identify(cmd, append(without(strings.Split(name, ":"), ""), rest...))
	}

	if cmds, err := cmd.Subcommands(); err != nil {
		return cmd, args, err
	} else if found := cmds.Find(name); found == nil {
		return nullCommand{parent: cmd, name: name}, rest, nil
	} else if subcmds, err := found.Subcommands(); err != nil {
		return found, rest, err
	} else if len(subcmds) > 0 {
		return identify(found, rest)
	} else {
		return found, rest, nil
	}
}

func isFlag(s string) bool {
	return strings.HasPrefix(s, "-")
}

func without(slice []string, exception string) (result []string) {
	for _, s := range slice {
		if s != exception {
			result = append(result, s)
		}
	}
	return
}
