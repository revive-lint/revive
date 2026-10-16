# Developer's Guide

This document explains how to build, test, and develop features for revive.

## Installation

Clone the project:

```bash
git clone git@github.com:revive-lint/revive.git
cd revive
```

## Build

In order to build the project run:

```bash
make build
```

The command will produce the `revive` binary in the root of the project.

## Test

Run the whole suite with:

```bash
make test
```

The tests come in two kinds:

- **Rule tests** in [`test/`](./test), one `_test.go` per rule, driven by the Go source fixtures in [`testdata/`](./testdata).
  They call the linter as a library and compare the reported failures with the `// MATCH` comments in the fixtures.
  Run a single rule with `go test -run TestUnusedParam ./test/...`.
- **End-to-end tests** in [`cli/e2e_test.go`](./cli/e2e_test.go) (`package cli_test`), run with `go test -run TestRevive ./cli/`.
  They execute revive as a separate process and assert on its exit code and on its output on stdout and stderr,
  which covers the behavior a library call cannot:
  flag parsing, `-set_exit_status`, the `warning-code` and `error-code` config settings, `-exclude`, and config file errors.

The end-to-end tests do not build a binary. The test binary re-executes itself with `REVIVE_E2E_RUN_REVIVE=1` set,
and its `TestMain` sees that variable and calls `RunRevive()` instead of running any test, so coverage still counts the CLI code.

Fixtures for these tests are in [`cli/testdata/`](./cli/testdata): `clean.go` reports nothing, `bad.go` reports two failures,
and each `.toml` file sets up one configuration under test. The subprocess runs in that directory with `HOME`, `USERPROFILE`,
and `XDG_CONFIG_HOME` pointed at temporary directories, so a `revive.toml` in your home directory cannot affect the outcome.
A case that needs one writes it there itself through the `homeConfig` and `xdgConfig` fields. That is how these tests cover
the default config discovery: `$XDG_CONFIG_HOME/revive.toml` first, then `$HOME/revive.toml`, then the built-in defaults.

## Logging

By default, any logging output is disabled when `REVIVE_LOG_LEVEL` is unset or empty.
You can enable and customize the log level using the `REVIVE_LOG_LEVEL` environment variable.
Supported values are:

- `debug`: log all messages including debug-level information
- `info`: log informational messages and above
- `warn`: log warnings and errors; also used as a fallback when `REVIVE_LOG_LEVEL` is set to an invalid value
- `error`: log errors only

Logs are output to stderr:

```sh
REVIVE_LOG_LEVEL=debug go run main.go
```

### Timing

At the `debug` level, revive also reports how long each part of a run takes. Use it to find out why a run is slow:

```sh
REVIVE_LOG_LEVEL=debug revive ./... 2>timing.log
```

Every measurement is a `phase completed` record carrying a `phase` and a `duration`:

| Phase | Reported | Measures |
| --- | --- | --- |
| `config-load` | once | reading and parsing the configuration file |
| `rules-setup` | once | instantiating and configuring the enabled rules |
| `package-resolution` | once | expanding the given patterns into packages and files |
| `go-version-detection` | once | locating `go.mod` files and detecting the Go version of every package |
| `package-read-parse` | per package | reading the files and parsing them into ASTs |
| `type-check` | per package | type checking, triggered lazily by the first typed rule |
| `scan-sortable` | per package | collecting the types implementing `sort.Interface` |
| `package-lint` | per package | the whole package: read, parse, type check, and all rules |
| `disabled-intervals` | per file | scanning the file for `revive:disable` directives |
| `file-lint` | per file | applying all rules to the file |
| `lint` | once | all packages |
| `format` | once | finishing the formatter once every failure has been delivered to it |
| `total` | once | the whole run |

The formatter runs concurrently with linting, so `format` only covers what is left after `lint`, not the whole formatting work.
Conversely, `file-lint` and `package-lint` include the time spent waiting for the formatter to accept each failure:
with a slow formatter or many failures, they can be high even though the rules are fast.

After linting, each rule gets a `rule timing` record with its total duration, the number of applications, and the average, slowest rule first.
A typed rule's duration includes the package type check when that rule triggers it, or waits for it while another file triggers it:
compare it with the `type-check` records before blaming the rule.

Timing is only collected when `REVIVE_LOG_LEVEL=debug`.

## Coding standards

Follow [the instructions](./.github/instructions/go.instructions.md) which contain Go coding standards and conventions used by both humans and
GitHub Copilot.

## Development of rules

If you want to develop a new rule, follow as an example the already existing rules in the
[rule package](https://github.com/revive-lint/revive/tree/master/rule) and check your change against
[the rule checklist](./.github/instructions/rule.instructions.md), which GitHub Copilot also uses when reviewing pull requests.

When adding a new rule that does not require type information (for example, a rule that does not call `file.Pkg.TypeCheck()`
and works purely on syntax/AST), add its name to `untyped.toml` and keep that file in sync with any such rules.

Each rule needs to implement the `lint.Rule` interface:

```go
type Rule interface {
	Name() string
	Apply(file *File, arguments Arguments) []Failure
}
```

All rules with a configuration must implement `lint.ConfigurableRule` interface:

```go
type ConfigurableRule interface {
	Configure(arguments Arguments) error
}
```

Add a compile-time assertion next to the rule type so that the interface is guaranteed to be implemented:

```golang
var _ lint.ConfigurableRule = (*ArgumentsLimitRule)(nil)
```

The `Arguments` type is an alias of the type `[]any`. The arguments of the rule are passed from the configuration file.

Every `lint.Failure` reported by a rule must set `Category` to one of the `lint.FailureCategory*` constants defined in
[`lint/failure.go`](/lint/failure.go). Failures without a category are rejected by the test harness.

Keep the failure messages of existing rules stable: do not reword one unless the change is necessary, for example because
the message became wrong after a behavior change. Users grep for these strings, match them in `// MATCH` test fixtures,
and baseline them in CI, so a cosmetic rewording breaks setups without making anything better.

### Example

Let's suppose we have developed a rule called `BanStructNameRule` which disallow us to name a structure with a given identifier.
We can set the banned identifier by using the TOML configuration file:

```toml
[rule.ban-struct-name]
arguments = ["Foo"]
```

With the snippet above we:

- Enable the rule with the name `ban-struct-name`. The `Name()` method of our rule should return a string that matches `ban-struct-name`.
- Configure the rule with the argument `Foo`.
The list of arguments will be passed to `Apply(*File, Arguments)` together with the target file we're linting currently.

A sample rule implementation can be [found here](/rule/argument_limit.go).

## Development of formatters

If you want to develop a new formatter, follow as an example the already existing formatters in the [formatter package](https://github.com/revive-lint/revive/tree/master/formatter).

All formatters should implement the following interface:

```go
type Formatter interface {
	Format(<-chan Failure, RulesConfig) (string, error)
	Name() string
}
```

## Lint

### Lint Markdown files

We use [markdownlint](https://github.com/DavidAnson/markdownlint),
[markdown-toc](https://github.com/jonschlinkert/markdown-toc),
and [mdsf](https://github.com/hougesen/mdsf) to check Markdown files.
`markdownlint` verifies document formatting, such as line length and empty lines.
`markdown-toc` checks the entries in the table of contents.
`mdsf` is responsible for formatting code snippets.

1. Install [markdownlint-cli2](https://github.com/DavidAnson/markdownlint-cli2#install).
2. Install [markdown-toc](https://github.com/jonschlinkert/markdown-toc#quick-start).
3. Install [mdsf](https://mdsf.mhouge.dk/#installation) and formatters:
    - [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) for `go`: `go install golang.org/x/tools/cmd/goimports@latest`
    - [shfmt](https://github.com/mvdan/sh#shfmt) for `sh, shell, bash`: `go install mvdan.cc/sh/v3/cmd/shfmt@latest`
    - [taplo](https://taplo.tamasfe.dev/cli/installation/binary.html) for `toml`
4. Run the following command to check formatting:

```shellsession
$ markdownlint-cli2 .
Finding: *.{md,markdown} *.md
Found:
 CODE_OF_CONDUCT.md
 CONTRIBUTING.md
 DEVELOPING.md
 README.md
 RULES_DESCRIPTIONS.md
Linting: 5 file(s)
Summary: 0 error(s)
```

_The `markdownlint-cli2` tool automatically uses the config file [.markdownlint-cli2.yaml](./.markdownlint-cli2.yaml)._
\
4. Run the following command to check TOC:

```sh
markdown-toc --maxdepth 4 --no-first1 --bullets "-" -i README.md && git diff --exit-code README.md
markdown-toc --maxdepth 2 --no-first1 --bullets "-" -i RULES_DESCRIPTIONS.md && git diff --exit-code RULES_DESCRIPTIONS.md
```

\
5. Run the following commands to verify and format code snippets:

```sh
mdsf verify .
```

```sh
mdsf format .
```

_Note: Use `golang` for Go code snippets that are intentionally non-compilable.
However, it is recommended to avoid this and use `go` whenever possible._

## Releasing

Releases are cut from `master` by pushing a tag; there are no release branches and no backports.
Pushing a tag triggers [`release.yml`](.github/workflows/release.yml), which runs GoReleaser and publishes the `ghcr.io/revive-lint/revive` image.

### When to release

- **Minor (`v1.N.0`)**: every 3 months, and never later than 6 months after the previous minor,
  so the supported Go version is bumped at least once per Go release cycle.
  Ship whatever is on `master`.
  When possible, tag 1–2 weeks before the next expected [golangci-lint](https://github.com/golangci/golangci-lint/releases) minor,
  so its dependency bump lands before their release.
- **Patch (`v1.N.M`)**: within 2 weeks of merging a fix for a regression, a panic, a false positive/negative in a rule enabled by default
  (in revive or golangci-lint), or a Go version compatibility issue.
  If `feat:` commits or new rules have already been merged since the last tag, bump the minor instead.
- **Major (`v2.0.0`)**: not time-driven; when the items tracked in [#1391](https://github.com/revive-lint/revive/issues/1391) are ready.

Go version: raise `go` in `go.mod` to the previous Go release in the first minor after a new Go version is released,
matching the [Go release policy](https://go.dev/doc/devel/release#policy).

### How to release

1. Check `git log <last-tag>..master`: any `feat`/`feature` commit (including scoped forms) or new rule means a minor, otherwise a patch.
2. From an up-to-date `master`, tag and push with `<remote>` set to the remote pointing at `revive-lint/revive`
   (`origin` in a direct clone, `upstream` when working from a fork — check with `git remote -v`):
   `git switch master && git pull --ff-only <remote> master && git tag vX.Y.Z && git push <remote> vX.Y.Z`.
3. Edit the auto-generated release notes into the sections used by previous releases, dropping any that have no entries:
   `## What's Changed`, with `### Features`, `### Bug Fixes`, `### Documentation`, `### Dependency Updates` and `### Chores & Internal`,
   followed by `## New Contributors` and the `**Full Changelog**` link.
   Every merged PR since the last tag belongs in exactly one section.
   Call out the Go version bump and any behavior changes that will produce new findings on existing code.

## Website

The documentation website <https://revive.run/> lives in a separate repository: [revive-lint/revive.run](https://github.com/revive-lint/revive.run).

Most of its content is generated from this repository, so there is no need to edit the website when changing docs here:

- `/docs` is generated from `README.md`
- `/r` is generated from `RULES_DESCRIPTIONS.md`
- `/images` is generated from `assets/`

The website is rebuilt from the latest revive release twice a month by a scheduled workflow,
so documentation changes appear on the website after the next release.

Keep the headings in `RULES_DESCRIPTIONS.md` stable: the `/r/#<rule>` links printed by revive rely on them.
