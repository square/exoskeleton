# Changelog

All notable changes to this project will be documented in this file. See [conventional commits](https://www.conventionalcommits.org/) for commit guidelines.

## [1.8.0](https://github.com/square/exoskeleton/compare/v1.7.0..v1.8.0)&nbsp;&nbsp;·&nbsp;&nbsp;2026-01-16

- **refactor:** Extract Contracts from discoverer and allow clients to supply their own (#48)&nbsp;&nbsp;·&nbsp;&nbsp;[b8f0e8b](https://github.com/square/exoskeleton/commit/b8f0e8b3a941792d9cd5ce35df0a3b09f5080bbf)

## [1.7.0](https://github.com/square/exoskeleton/compare/v1.6.6..v1.7.0)&nbsp;&nbsp;·&nbsp;&nbsp;2025-12-12

- **perf:** Parallelize operations when we're evaluating subcommands' contracts (#46)&nbsp;&nbsp;·&nbsp;&nbsp;[b68e0ab](https://github.com/square/exoskeleton/commit/b68e0abfdcc46a0369d1885859c9e43ccc6bd60a)
- **refactor:** Extract a parameterizeable `Expand` from `Flatten` and `expandModules` (#45)&nbsp;&nbsp;·&nbsp;&nbsp;[7792b0d](https://github.com/square/exoskeleton/commit/7792b0d36acf1ecdadfca82314ead235e994f0fc)
- **chore:** Add a workflow that maintains a CHANGELOG.md (#40)&nbsp;&nbsp;·&nbsp;&nbsp;[d34fbef](https://github.com/square/exoskeleton/commit/d34fbefd07d6b2ab95cf7da6b6280a34dfa1db55)

## [1.6.6](https://github.com/square/exoskeleton/compare/v1.6.5..v1.6.6)&nbsp;&nbsp;·&nbsp;&nbsp;2024-09-27

- **feat:** Add `WithExecutor` option that allows clients to manipulate `*exec.Cmd` before it is run (#39)&nbsp;&nbsp;·&nbsp;&nbsp;[afef7f4](https://github.com/square/exoskeleton/commit/afef7f4412e10575a8d86c2e979b3e0983be0078)

## [1.6.5](https://github.com/square/exoskeleton/compare/v1.6.4..v1.6.5)&nbsp;&nbsp;·&nbsp;&nbsp;2024-09-13

- **fix:** Don't suggest the same completion more than once (#38)&nbsp;&nbsp;·&nbsp;&nbsp;[c0ccebe](https://github.com/square/exoskeleton/commit/c0ccebec92f9ba786e701f9a6114aec2145e8729)
- **fix:** Do not allow a command to be renamed after it was discovered (#37)&nbsp;&nbsp;·&nbsp;&nbsp;[1f4e1ba](https://github.com/square/exoskeleton/commit/1f4e1ba35c0a78e400d43f6144c25200b030093d)

## [1.6.4](https://github.com/square/exoskeleton/compare/v1.6.3..v1.6.4)&nbsp;&nbsp;·&nbsp;&nbsp;2024-09-10

- **fix:** Allow `--describe-commands` to provide a an empty summary (#36)&nbsp;&nbsp;·&nbsp;&nbsp;[9715947](https://github.com/square/exoskeleton/commit/9715947282ea5ba992b784c0c96dda8a3db5e89d)
- **chore:** Improve error messages (#35)&nbsp;&nbsp;·&nbsp;&nbsp;[1a8fa9e](https://github.com/square/exoskeleton/commit/1a8fa9e0adf4b00331a999e74580ab7d0d07d96b)

## [1.6.3](https://github.com/square/exoskeleton/compare/v1.6.2..v1.6.3)&nbsp;&nbsp;·&nbsp;&nbsp;2024-09-06

- **fix:** Only foreground subprocesses if we have a controlling terminal (#34)&nbsp;&nbsp;·&nbsp;&nbsp;[7e39500](https://github.com/square/exoskeleton/commit/7e39500b58cd245a82f51afcfbe8c1dfc6428793)

## [1.6.2](https://github.com/square/exoskeleton/compare/v1.6.1..v1.6.2)&nbsp;&nbsp;·&nbsp;&nbsp;2024-08-23

- **fix:** Don't foreground subprocesses if they return ENOTTY (#33)&nbsp;&nbsp;·&nbsp;&nbsp;[051267c](https://github.com/square/exoskeleton/commit/051267cb388d9547dcc6ca3772da560b05962c2c)

## [1.6.1](https://github.com/square/exoskeleton/compare/v1.6.0..v1.6.1)&nbsp;&nbsp;·&nbsp;&nbsp;2024-08-21

- **fix:** Foreground subcommands so that signals get relayed to them instead of to the exoskeleton (#31)&nbsp;&nbsp;·&nbsp;&nbsp;[fb62f1e](https://github.com/square/exoskeleton/commit/fb62f1ebfad83737389c466496d056a61fee7304)

## [1.6.0](https://github.com/square/exoskeleton/compare/v1.5.0..v1.6.0)&nbsp;&nbsp;·&nbsp;&nbsp;2024-08-20

- **feat:** Render menus using Go templates and allow clients to supply their own (#27)&nbsp;&nbsp;·&nbsp;&nbsp;[023dd9d](https://github.com/square/exoskeleton/commit/023dd9d3141227b3ea127ec6f52731d316b1dcdb)
- **feat:** Allow registering a callback to be invoked whenever a command is successfully identified (#30)&nbsp;&nbsp;·&nbsp;&nbsp;[7107e16](https://github.com/square/exoskeleton/commit/7107e169a032d7e6e45e41899d29089f01eb6689)
- **fix:** Do not normalize `--help` or `-h` when they don't immediately follow an identifiable command (#26)&nbsp;&nbsp;·&nbsp;&nbsp;[1b143ce](https://github.com/square/exoskeleton/commit/1b143cec40fa1b23bd0074e780bf63817fd53b5d)
- **refactor:** Move `symlink-test` to `fixtures/edge-cases`&nbsp;&nbsp;·&nbsp;&nbsp;[05288c8](https://github.com/square/exoskeleton/commit/05288c821d1aa2c2818385590c2ddf46fd8e2a1d)
- **chore:** Stop wrapping `Exec`'s errors (#21)&nbsp;&nbsp;·&nbsp;&nbsp;[11d135e](https://github.com/square/exoskeleton/commit/11d135ed495ccaa33655de891bc3a1cd01439b62)
- **chore:** Remove `summaryCache` (defer that to clients) (#29)&nbsp;&nbsp;·&nbsp;&nbsp;[7d04a9d](https://github.com/square/exoskeleton/commit/7d04a9d2045ec3a465c4096856c02789ea14d8f6)

## [1.5.0](https://github.com/square/exoskeleton/compare/v1.4.1..v1.5.0)&nbsp;&nbsp;·&nbsp;&nbsp;2024-06-11

- **feat:** Allow `help`, `which`, and `complete` to be overridden by prepended commands (#24)&nbsp;&nbsp;·&nbsp;&nbsp;[4937a57](https://github.com/square/exoskeleton/commit/4937a57b04a57515c192df26b96e852e007abf46)
- **refactor:** Export HelpExec, CompleteExec, and WhichExec, the implementations of `help`, `complete`, and `which` (#25)&nbsp;&nbsp;·&nbsp;&nbsp;[db0801f](https://github.com/square/exoskeleton/commit/db0801f18163ea2059f2d0c8181f463b6e5026b0)

## [1.4.1](https://github.com/square/exoskeleton/compare/v1.4.0..v1.4.1)&nbsp;&nbsp;·&nbsp;&nbsp;2024-06-11

- **feat:** Add `SetName` option for setting the name of the entrypoint (#23)&nbsp;&nbsp;·&nbsp;&nbsp;[05e749b](https://github.com/square/exoskeleton/commit/05e749b00c33ee494c4033fc725d54a5db7d6241)

## [1.4.0](https://github.com/square/exoskeleton/compare/v1.3.1..v1.4.0)&nbsp;&nbsp;·&nbsp;&nbsp;2024-06-06

- **fix:** Evaluate `--describe-commands` lazily (#15)&nbsp;&nbsp;·&nbsp;&nbsp;[e3ae428](https://github.com/square/exoskeleton/commit/e3ae4281804702101c6d5f4aa8281ed507e6e3b1)
- **fix:** Categorize bad responses from `--describe-commands` as failures to implement Command's API not Discovery errors (#20)&nbsp;&nbsp;·&nbsp;&nbsp;[b5cd2ce](https://github.com/square/exoskeleton/commit/b5cd2ce843e8fdd629200abd591c950fca3267f5)
- **fix** [**breaking**]**:** Allow Module's `Subcommands()` and Entrypoint's `Identify()` to return CommandErrors (#19)&nbsp;&nbsp;·&nbsp;&nbsp;[5559f4c](https://github.com/square/exoskeleton/commit/5559f4cb5b67d21a45fd7ca2986de41aa7aea6df)
- **refactor:** Move code from `executable_command.go` and `directory_module.go` to `contract.go` (#17)&nbsp;&nbsp;·&nbsp;&nbsp;[bfe2bab](https://github.com/square/exoskeleton/commit/bfe2babb39f23756488a228281163d474dd0039e)

## [1.3.1](https://github.com/square/exoskeleton/compare/v1.3.0..v1.3.1)&nbsp;&nbsp;·&nbsp;&nbsp;2024-05-22

- **fix:** Return `[]` from `discoverIn` instead of `nil` (#14)&nbsp;&nbsp;·&nbsp;&nbsp;[d66edec](https://github.com/square/exoskeleton/commit/d66edec2c9b812e82cdfff350a27d5bd01078e30)

## [1.3.0](https://github.com/square/exoskeleton/compare/v1.2.0..v1.3.0)&nbsp;&nbsp;·&nbsp;&nbsp;2024-05-01

- **feat:** Allow passing arbitrary arguments through to executables along with `--help` (#13)&nbsp;&nbsp;·&nbsp;&nbsp;[0ab680c](https://github.com/square/exoskeleton/commit/0ab680cf1adfec713448e2e803fe7c31f804a7cc)

## [1.2.0](https://github.com/square/exoskeleton/compare/v1.1.3..v1.2.0)&nbsp;&nbsp;·&nbsp;&nbsp;2024-04-25

- **feat:** Define a contract that allows a single executable to represent a subtree of commands (#12)&nbsp;&nbsp;·&nbsp;&nbsp;[cc1b1f2](https://github.com/square/exoskeleton/commit/cc1b1f25a83e1f47be98d096b08bb1ca35656c45)
- **refactor:** Move the `summaryCache` out of `menu.go`&nbsp;&nbsp;·&nbsp;&nbsp;[c54fc41](https://github.com/square/exoskeleton/commit/c54fc41f5a032a0ff6783bf47e452518fd550faa)
- **refactor:** Make `executableCommand` the receiver of `getMessageFromExecution`&nbsp;&nbsp;·&nbsp;&nbsp;[b487899](https://github.com/square/exoskeleton/commit/b4878995ccd2b9a54f3c928bac3653494d783e37)
- **refactor:** Extract `Command` from `executable`'s functions&nbsp;&nbsp;·&nbsp;&nbsp;[494c98e](https://github.com/square/exoskeleton/commit/494c98e7cfdf77ef2bd38a8bda561fc31cb5272f)
- **refactor:** Embed `discoverer` into `directoryModule` instead of using a reference&nbsp;&nbsp;·&nbsp;&nbsp;[ef269f3](https://github.com/square/exoskeleton/commit/ef269f3cef8af88c568a0b36af83f39577b8480a)
- **refactor:** Replace `path.Join` with `filepath.Join`&nbsp;&nbsp;·&nbsp;&nbsp;[19cb0a0](https://github.com/square/exoskeleton/commit/19cb0a06cb65708df0dce9306a9858931a0aa665)
- **refactor:** Rename `p` to `path`&nbsp;&nbsp;·&nbsp;&nbsp;[ac43282](https://github.com/square/exoskeleton/commit/ac43282d4d5df44d6ec10272179a87caa3b5283b)
- **refactor:** Extract `buildCommand` from `discoverIn` and test it&nbsp;&nbsp;·&nbsp;&nbsp;[9caec71](https://github.com/square/exoskeleton/commit/9caec719d06dc09b98b517a9fcaa32f76553b019)

## [1.1.3](https://github.com/square/exoskeleton/compare/v1.1.2..v1.1.3)&nbsp;&nbsp;·&nbsp;&nbsp;2024-04-09

- **fix:** Trim all trailing newline characters instead of just the first (#11)&nbsp;&nbsp;·&nbsp;&nbsp;[cece69e](https://github.com/square/exoskeleton/commit/cece69eab8be258616beff0465918a4865412253)

## [1.1.2](https://github.com/square/exoskeleton/compare/v1.1.1..v1.1.2)&nbsp;&nbsp;·&nbsp;&nbsp;2024-04-09

- **fix:** Capture only standard output from invocations using `--summary` and `--help` (#10)&nbsp;&nbsp;·&nbsp;&nbsp;[28bb72a](https://github.com/square/exoskeleton/commit/28bb72ab95fe781d22bb1b604d983185cd45e699)

## [1.1.1](https://github.com/square/exoskeleton/compare/v1.1.0..v1.1.1)&nbsp;&nbsp;·&nbsp;&nbsp;2024-04-03

- **feat:** Add `IsEmbedded` to identify embedded commands (#7)&nbsp;&nbsp;·&nbsp;&nbsp;[58a46da](https://github.com/square/exoskeleton/commit/58a46dac2a5550c9e2946e52b619ab6995737748)
- **feat:** Add `CompleteFiles` to provide filesystem completions for built-in commands (#8)&nbsp;&nbsp;·&nbsp;&nbsp;[2a45844](https://github.com/square/exoskeleton/commit/2a45844073d63f27327a52bcc33e082f8a509486)
- **feat:** Add `IsModule` to identify modules (#9)&nbsp;&nbsp;·&nbsp;&nbsp;[595445a](https://github.com/square/exoskeleton/commit/595445a48f8c0ed03843b3593f8b23d7686ffddb)

## [1.1.0](https://github.com/square/exoskeleton/compare/v1.0.2..v1.1.0)&nbsp;&nbsp;·&nbsp;&nbsp;2024-04-01

- **feat** [**breaking**]**:** Allow defining embedded modules in addition to embedded commands (#3)&nbsp;&nbsp;·&nbsp;&nbsp;[0a13a0b](https://github.com/square/exoskeleton/commit/0a13a0b666dca1b7b433eba050fa231cc6c74c86)
- **chore:** Remove the `--with-modules` flag from `help` (#2)&nbsp;&nbsp;·&nbsp;&nbsp;[4170ba9](https://github.com/square/exoskeleton/commit/4170ba91a264c6535fa25926cd084a0f1d18f1c4)

## [1.0.2](https://github.com/square/exoskeleton/compare/v1.0.1..v1.0.2)&nbsp;&nbsp;·&nbsp;&nbsp;2024-03-14


## [1.0.1]&nbsp;&nbsp;·&nbsp;&nbsp;2024-03-08

- **feat:** Open-source exoskeleton&nbsp;&nbsp;·&nbsp;&nbsp;[af64657](https://github.com/square/exoskeleton/commit/af64657ad6bb7f1eaebada5c1519df31b3b03fcb)

<!-- generated by git-cliff -->
