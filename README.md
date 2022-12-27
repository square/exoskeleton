# Exoskeleton

Exoskeleton is a library for creating modern multi-CLI applications whose commands are external to them.

[![](https://img.shields.io/github/actions/workflow/status/square/exoskeleton/ci.yml?branch=main&longCache=true&label=CI&logo=github%20actions&logoColor=fff)](https://github.com/square/exoskeleton/actions?query=workflow%3ACI)
[![Go Reference](https://pkg.go.dev/badge/github.com/square/exoskeleton.svg)](https://pkg.go.dev/github.com/square/exoskeleton)

## Why Exoskeleton?

Exoskeleton is similar to frameworks like [Cobra](cobra) and [Oclif](oclif) in that it simplifies building subcommand-based command line applications (like `git clone`, `git log`) by providing:
- nested subcommands and automatically generated menus of subcommands
- tab-completion for shells
- suggestions for mistyped commands ("Did you mean ...?")

Exoskeleton differs from other CLI frameworks in that each subcommand maps to a separate, standalone executable. This allows subcommands to be implemented in different languages and released on different schedules!

Exoskeleton's goal is to create a common entrypoint, better discoverability, and a cohesive experience for a decentralized suite of commandline tools.

## Creating an Exoskeleton

Install Exoskeleton by running
```sh
go get -u github.com/square/exoskeleton@latest
```

The simplest Exoskeleton application looks like this:
```golang
package main

import (
	"os"

	"github.com/square/exit"
	"github.com/square/exoskeleton"
)

func main() {
	paths := []string{
		// paths
		// you want Exoskeleton to search
		// for subcommands
	}

	// 1. Create a new Command Line application that looks for subcommands in the given paths.
	cli, _ := exoskeleton.New(paths)

	// 2. Identify the subcommand being invoked from the arguments.
	cmd, args := cli.Identify(os.Args[1:])

	// 3. Execute the subcommand.
	err := cmd.Exec(cli, args, os.Environ())

	// 4. Exit the program with the exit code the subcommand returned.
	os.Exit(exit.FromError(err))
}
```

① An exoskeleton is constructed with an array of paths to search for [subcommands](subcommands). ② It uses the arguments to identify which subcommand is being invoked. ③ The subcommand is executed ④ and the program exits (`0` if `err` is `nil`, `1` or a [semantic exit code][exit] otherwise).

To see this in action, take a look at the [the Hello World example project](hello_world).

In the real world, a CLI might also:
1. Customize the exoskeleton by passing [options](options) to `exoskeleton.New`
2. Add business logic between ② `Identify` and ③ `Exec` or ③ `Exec` and ④ `os.Exit`

> [!TIP]
> At Square, we use the [OnCommandNotFound](OnCommandNotFound) callback to install subcommands on-demand, check for updates after constructing the exoskeleton, and wrap `Exec` to emit usage metrics.

# Subcommands

## Creating Subcommands

Subcommands can be either shell scripts or binaries.
1. They MUST be executable.
1. They SHOULD output help text to be displayed when the user invokes them with `--help`.
1. They MAY respond to `--summary` by outputting a summary of their purpose to be displayed in a menu of commands.
1. They MAY respond to `--complete <input>` by outputting a list of shell-completions for `<input>`.

### Help and Summary Text

Compiled binaries should parse their arguments for the `--help` and `--summary` flags. They should do this early in execution before any expensive set up, write the text to standard out, and exit successfully.

Shell scripts which start with the shebang (`#!`) may respond to `--help` and `--summary` flags or may choose to document themselves with magic comments (`# HELP: <help text follows>`, `# SUMMARY: <summary line follows>`). See [examples/hello_world/libexec/ls](ls) and [examples/hello_world/libexec/rm](rm) for examples.

### Completions

Exoskeleton uses [shellcomp](shellcomp) (the API that Cobra developed) to separate shell-specific logic for implementing completions from the logic for producing the suggestions themselves.

Shellcomp scripts invoke exoskeleton with `--complete <WHATEVER THE USER TYPED>` and expect to receive a list of suggestions on standard output.

> [!NOTE]
> Imagine `git` is implemented with Exoskeleton.
>
> If the user types
> ```
> $ git che<tab>
> ```
> then the shellcomp scripts will execute:
> ```
> $ git complete che
> ```
> and Exoskeleton will suggest completions from the list of commands it knows and output:
> ```
> checkout
> :4
> ```
> and the shellcomp scripts will complete the command `git checkout`.
>
> If the user types
> ```
> $ git checkout lail/<tab>
> ```
> then the shellcomp scripts will execute:
> ```
> $ git complete checkout lail/
> ```
> and Exoskeleton will dispatch the completion to the `checkout` command, executing this:
> ```
> $ $(git which checkout) --complete -- checkout lail/
> ```
> and, at this point, if `git checkout` handles `--complete`, it may list branches that start with `lail/`.

See [shellcomp's docs](shellcomp) for implementing completions for a subcommand.

Call [exoskeleton.GenerateCompletionScript](GenerateCompletionScript) to generate the shellcomp scripts for your project.

> [!TIP]
> At Square, we call this function in our `Makefile` and distribute artifacts for Bash and Zsh with releases.

## Menus

In Exoskeleton, each subcommand maps to a separate, standalone executable. _Submenus_ map to directories. The directory must contain a file named `.exoskeleton` (this is [configurable](WithModuleMetadataFilename)).

Take a look at the `dir` module in [the Hello World example project](hello_world).


[cobra]: https://github.com/spf13/cobra
[exit]: https://github.com/square/exit
[hello_world]: https://github.com/square/exoskeleton/tree/main/examples/hello_world
[ls]: https://github.com/square/exoskeleton/tree/main/examples/hello_world/libexec/ls
[oclif]: https://oclif.io/
[options]: https://pkg.go.dev/github.com/square/exoskeleton#Option
[rm]: https://github.com/square/exoskeleton/tree/main/examples/hello_world/libexec/rm
[sub]: https://github.com/qrush/sub
[subcommands]: #subcommands
[GenerateCompletionScript]: https://pkg.go.dev/github.com/square/exoskeleton#GenerateCompletionScript
[OnCommandNotFound]: https://pkg.go.dev/github.com/square/exoskeleton#OnCommandNotFound
[WithModuleMetadataFilename]: https://pkg.go.dev/github.com/square/exoskeleton#WithModuleMetadataFilename
[shellcomp]: https://github.com/square/exoskeleton/tree/main/pkg/shellcomp
