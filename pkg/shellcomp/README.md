# shellcomp

A commandline tool can provide completions by responding to the `--complete` flag and making suggestions based on any arguments that follow the flag. The tool should write a list of suggestions to standard output (separated by newlines) and a [directive](#directives)

Imaging that typing `git che<tab>` causes the shell to complete the word `checkout` and then typing `git checkout <tab>` causes the shell to list local branches.

If `git` implemented its completion scripts using this API, it would be expected to interact like this:
```
$ git --complete che
checkout
:4

$ git --complete checkout ""
<list of local branches>
:4
```

The directive `:4` tells the completion script not to suggest files from the working directory.

### Directives

The directive is a bit flag that can combine the following options:

| Directive | Name | Use for |
| --------: | :-- | :-- |
|         0 | `DirectiveDefault` | The shell should perform its default behavior after providing the returned completions |
|         1 | `DirectiveError` | The shell should ignore completions: an error occurred |
|         2 | `DirectiveNoSpace` | The shell should not add a space after completing a word (useful for completing arguments to flags that end with `=`) |
|         4 | `DirectiveNoFileComp` | The shell should not suggest files from the working directory |
|         8 | `DirectiveFilterFileExt` | The shell should suggest files from the working directory and use the returned completions as file extension filters instead of suggestions |
|         16 | `DirectiveFilterDirs` | The shell should suggest directories and use the returned completions to identify the directory in which to search |


### Go

Go projects may import the package `"github.com/square/exoskeleton/pkg/shellcomp"` and call `os.Stdout.Write(shellcomp.Marshal(completions, directive, false))`

