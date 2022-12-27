// Package exoskeleton allows you to create modern multi-CLI applications
// whose subcommands are external to them.
//
// You use it by defining an Entrypoint and a list of paths where external
// subcommands may be found.
//
// Here's an example:
//
//	Assuming we have two executables that respond to `--summary` like this:
//
//	   $ ~/libexec/rm --summary
//	   Remove directory entries
//
//	   $ ~/libexec/ls --summary
//	   List directory contents
//
//	And a binary, `example`, that is implemented like this:
//
//	   package main
//
//	   import (
//	       "os"
//
//	       "github.com/squareup/exoskeleton"
//	   )
//
//	   func main() {
//	       exoskeleton.Exec([]string{os.Getenv("HOME") + "/libexec"})
//	   }
//
//	Then our example program will behave like this:
//
//	   $ ./example
//	   USAGE
//	      example <command> [<args>]
//
//	   COMMANDS
//	      ls  List directory contents
//	      rm  Remove directory entries
//
//	   Run example help <command> to print information on a specific command.
//
//	And running `example ls` will forward execution to `~/libexec/ls`.
package exoskeleton
