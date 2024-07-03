package exoskeleton

import "strings"

// Identify identifies the command being invoked.
//
// The function also returns any arguments that were not used to identify the command.
// For example, if the Go CLI were implemented with Exoskeleton, and we ran it
// with the rawArgs 'get -u github.com/square/exoskeleton', Identify would return
// the Get command and the arguments {'-u', 'github.com/square/exoskeleton'}.
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

	if IsNull(cmd) {
		e.commandNotFound(cmd)
	}

	// Recognize `--help` and `-h` as aliases for the built-in `help` command
	// only when they immediately follow an identifiable command.
	if !IsNull(cmd) && len(rest) > 0 && (rest[0] == "--help" || rest[0] == "-h") {
		return e.Identify(append(append([]string{"help"}, argsRelativeTo(cmd, e)...), rest[1:]...))
	}

	return cmd, rest, err
}

// identify uses args to identify a Command and returns the command and the rest
// of the commandline arguments or else {nil, args} if no Command is identified.
//
// Returns a CommandError if the command does not fulfill the contract
// for providing its subcommands.
func identify(m Module, args []string) (Command, []string, error) {
	if len(args) == 0 || isFlag(args[0]) {
		return m, args, nil
	}

	// Rewrite args like {"module:subcommand", "--flag"} to {"module", "subcommand", "--flag"}.
	// Do this just-in-time, non-destructively, while we're working on identifying a command.
	name, rest := args[0], args[1:]
	if strings.Contains(name, ":") {
		return identify(m, append(without(strings.Split(name, ":"), ""), rest...))
	}

	if cmds, err := m.Subcommands(); err != nil {
		return m, args, err
	} else if cmd := cmds.Find(name); cmd == nil {
		return nullCommand{parent: m, name: name}, rest, nil
	} else if submodule, ok := cmd.(Module); ok {
		return identify(submodule, rest)
	} else {
		return cmd, rest, nil
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
