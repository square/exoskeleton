package exoskeleton

import "strings"

// Identify identifies the command being invoked.
//
// The function also returns any arguments that were not used to identify the command.
// For example, if the Go CLI were implemented with Exoskeleton, and we ran it
// with the rawArgs 'get -u github.com/squarup/exoskeleton', Identify would return
// the Get command and the arguments {'-u', 'github.com/squarup/exoskeleton'}.
//
// If no command is identified, Identify invokes CommandNotFound callbacks and
// returns NullCommand.
func (e *Entrypoint) Identify(rawArgs []string) (Command, []string, error) {
	cmd, args, err := identify(e, normalizeArgs(rawArgs))

	if n, ok := cmd.(nullCommand); ok {
		e.commandNotFound(n)
	}

	return cmd, args, err
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

// normalizeArgs rewrites `<command> -h` and `<command> --help` to `help <command>`.
// It also treats `--complete <text>` as a synonym for `complete <text>` when the flag
// is the first argument.
func normalizeArgs(args []string) []string {
	normalizedArgs := []string{}

	for i, arg := range args {
		if arg == "--" {
			normalizedArgs = append(normalizedArgs, args[i:]...)
			break

		} else if arg == "-h" || arg == "--help" {
			normalizedArgs = append([]string{"help"}, normalizedArgs...)

		} else if arg == "--complete" && i == 0 {
			normalizedArgs = append([]string{"complete"}, normalizedArgs...)

		} else {
			normalizedArgs = append(normalizedArgs, arg)
		}
	}

	return normalizedArgs
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
