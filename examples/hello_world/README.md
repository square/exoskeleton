# Hello World

This project demonstrates a basic commandline application created with Exoskeleton.

[hello.go](hello.go) constructs a CLI that will find its subcommands in the path `./libexec`.

[libexec/dir/ls](ls) and [libexec/dir/rm](rm) are two Bash scripts that [can be picked up as subcommands](contract).

You can run the project like this:
```sh
git clone org-49461806@github.com:square/exoskeleton.git
cd exoskeleton/examples/hello_world
go build
./hello
```

`hello` displays the menu of subcommands:
```
USAGE
   hello <command> [<args>]

COMMANDS
   dir:  Utilities for working with directories

Run hello help <command> to print information on a specific command.
```
(`hello help` and `hello -h`, and `hello --help` do the same thing.)

`hello dir` lists the commands in that directory:
```
USAGE
   hello dir <command> [<args>]

COMMANDS
   ls  List directory contents
   rm  Remove directory entries

Run hello help dir <command> to print information on a specific command.
```

`hello help dir ls` describes the `ls` subcommand:
```
USAGE
   ls [FILE]...

DESCRIPTION
   For each operand that names a file of a type other than directory,
   ls displays its name as well as any requested, associated information.

   For each operand that names a file of type directory, ls displays
   the names of files contained within that directory, as well as any
   requested, associated information.
```
(As does `hello dir ls -h` and `hello dir ls --help`)

And `hello dir ls` runs [libexec/dir/ls](ls).


[hello.go]: https://github.com/square/exoskeleton/tree/main/examples/hello_world/hello.go
[ls]: https://github.com/square/exoskeleton/tree/main/examples/hello_world/libexec/dir/ls
[rm]: https://github.com/square/exoskeleton/tree/main/examples/hello_world/libexec/dir/rm
[contract]: https://github.com/square/exoskeleton/tree/main/README.md#contract
