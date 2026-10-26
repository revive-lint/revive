# revive

[![Build Status](https://github.com/revive-lint/revive/actions/workflows/test.yaml/badge.svg)](https://github.com/revive-lint/revive/actions/workflows/test.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mgechev/revive.svg)](https://pkg.go.dev/github.com/mgechev/revive)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/revive-lint/revive)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/15967/badge)](https://www.bestpractices.dev/projects/15967)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/revive-lint/revive/badge)](https://scorecard.dev/viewer/?uri=github.com/revive-lint/revive)

Fast, configurable, extensible, flexible, and beautiful linter for Go. Drop-in replacement of golint.
**`Revive` provides a framework for development of custom rules,
and lets you define a strict preset for enhancing your development & code review processes**.

<p align="center">
  <img src="./assets/logo.png" alt="" width="300">
  <br>
  Logo by <a href="https://github.com/hawkgs">Georgi Serev</a>
</p>

Here's how `revive` is different from `golint`:

- Allows to enable or disable rules using a configuration file.
- Allows to configure the linting rules with a TOML file.
- 2x faster running the same rules as golint.
- Provides functionality for disabling a specific rule or the entire linter for a file or a range of lines.
  - `golint` allows this only for generated files.
- Optional type checking. Most rules in golint do not require type checking.
If you disable them in the config file, revive will run over 6x faster than golint.
- Provides multiple formatters which let us customize the output.
- Allows to customize the return code for the entire linter or based on the failure of only some rules.
- _Everyone can extend it easily with custom rules or formatters._
- `Revive` provides more rules compared to `golint`.

<p align="center">
  <img src="./assets/demo.svg" alt="" width="700">
</p>

<!-- toc -->

- [Installation](#installation)
  - [Homebrew](#homebrew)
  - [Install from Sources](#install-from-sources)
  - [Docker](#docker)
  - [Manual Binary Download](#manual-binary-download)
- [Usage](#usage)
  - [Text Editors](#text-editors)
  - [GitHub Actions](#github-actions)
  - [Continuous Integration](#continuous-integration)
  - [Linter aggregators](#linter-aggregators)
    - [golangci-lint](#golangci-lint)
    - [MegaLinter](#megalinter)
  - [Command Line Flags](#command-line-flags)
  - [Sample Invocations](#sample-invocations)
  - [Comment Directives](#comment-directives)
  - [Configuration](#configuration)
  - [Custom Configuration](#custom-configuration)
  - [Recommended Configuration](#recommended-configuration)
  - [Rule-level file excludes](#rule-level-file-excludes)
- [Available Rules](#available-rules)
- [Configurable rules](#configurable-rules)
- [Available Formatters](#available-formatters)
  - [Friendly](#friendly)
  - [Stylish](#stylish)
  - [Default](#default)
  - [Plain](#plain)
  - [Unix](#unix)
  - [JSON](#json)
  - [NDJSON](#ndjson)
  - [Checkstyle](#checkstyle)
  - [SARIF](#sarif)
- [Extensibility](#extensibility)
  - [Writing a Custom Rule](#writing-a-custom-rule)
    - [Using `revive` as a library](#using-revive-as-a-library)
  - [Custom Formatter](#custom-formatter)
- [Speed Comparison](#speed-comparison)
  - [golint](#golint)
  - [revive's speed](#revives-speed)
- [Overriding colorization detection](#overriding-colorization-detection)
- [Who uses Revive](#who-uses-revive)
- [Contributors](#contributors)
  - [Maintainers](#maintainers)
  - [All](#all)
- [License](#license)

<!-- tocstop -->

## Installation

`revive` is available inside the majority of package managers.

[<img alt="Repology packaging status" src="https://repology.org/badge/vertical-allrepos/revive.svg">](https://repology.org/project/revive/versions)

### Homebrew

Install `revive` using [brew](https://brew.sh/):

```bash
brew install revive
```

To upgrade to the latest version:

```bash
brew upgrade revive
```

### Install from Sources

Install the latest stable release directly from source:

```bash
go install github.com/mgechev/revive@latest
```

To install the latest commit from the main branch:

```bash
go install github.com/mgechev/revive@HEAD
```

### Docker

You can run `revive` using Docker to avoid installing it directly on your system:

```bash
docker run -v "$(pwd)":/var/YOUR_REPOSITORY ghcr.io/mgechev/revive:v1.16.0 -config /var/YOUR_REPOSITORY/revive.toml -formatter stylish ./var/YOUR_REPOSITORY/...
```

_Note_: Replace `YOUR_REPOSITORY` with the path to your repository.

A volume must be mounted to share the current repository with the container.
For more details, refer to the [bind mounts Docker documentation](https://docs.docker.com/storage/bind-mounts/).

- `-v`: Mounts the current directory (`$(pwd)`) to `/var/YOUR_REPOSITORY` inside the container.
- `ghcr.io/mgechev/revive:v1.16.0`: Specifies the Docker image and its version.
- `revive`: The command to run inside the container.
- Flags like `-config` and `-formatter` are the same as when using the binary directly.

_Note_: images up to and including `v1.16.0` are published at `ghcr.io/mgechev/revive`.
Later releases are published at `ghcr.io/revive-lint/revive`.

### Manual Binary Download

Download the precompiled binary from the [Releases page](https://github.com/revive-lint/revive/releases):

1. Select the appropriate binary for your OS and architecture.
2. Extract the binary and move it to a directory in your `PATH` (e.g., `/usr/local/bin`).
3. Verify installation:

```bash
revive -version
```

## Usage

Since the default behavior of `revive` is compatible with `golint`, without providing any additional flags,
the only difference you'd notice is faster execution.

`revive` supports a `-config` flag whose value should correspond to a TOML file describing which rules to use for `revive`'s linting.
If not provided, `revive` will try to use a global config file (assumed to be located at `$HOME/revive.toml`).
Otherwise, if no configuration TOML file is found then `revive` uses a built-in set of default linting rules.

### Text Editors

- Support for VSCode via [vscode-go](https://code.visualstudio.com/docs/languages/go#_build-and-diagnose) by changing the `go.lintTool` setting to `revive`:

```json
{
  "go.lintTool": "revive",
}
```

- Support for GoLand via [File Watchers](https://dev.to/s0xzwasd/configure-revive-go-linter-in-goland-2ggl).
- Support for vim via [dense-analysis/ale](https://github.com/dense-analysis/ale).

  ```vim
  let g:ale_linters = {
  \   'go': ['revive'],
  \}
  ```

### GitHub Actions

- [Revive Action](https://github.com/marketplace/actions/revive-action) with annotation support

### Continuous Integration

[Codeac.io](https://www.codeac.io?ref=revive) - Automated code review service integrates with GitHub,
Bitbucket and GitLab (even self-hosted) and helps you fight technical debt.
Check your [pull-requests](https://www.codeac.io/documentation/pull-requests.html?ref=revive) with
[revive](https://www.codeac.io/documentation/revive-configuration.html?ref=revive) automatically.
(Free for open-source projects)

### Linter aggregators

#### golangci-lint

To enable `revive` in `golangci-lint` you need to add `revive` to the list of enabled linters:

```yaml
# golangci-lint configuration file
version: "2"
linters:
   enable:
     - revive
```

Then `revive` can be configured by adding an entry to the `linters.settings` section of the configuration, for example:

```yaml
# golangci-lint configuration file
linters:
  settings:
    revive:
      severity: warning
      rules:
        - name: atomic
        - name: line-length-limit
          severity: error
          arguments: [80]
        - name: unhandled-error
          arguments: ["fmt.Printf", "myFunction"]
```

The above configuration enables three rules of `revive`: _atomic_, _line-length-limit_ and _unhandled-error_ and passes some arguments to the last two.
The [Configuration](#configuration) section of this document provides details on how to configure `revive`.
Note that while `revive` configuration is in TOML, that of `golangci-lint` is in YAML or JSON.
See the [golangci-lint website](https://golangci-lint.run/usage/linters/#revive) for more information about configuring `revive`.

Please notice that if no particular configuration is provided, `revive` will behave as `golint` does, i.e. all `golint` rules are enabled
(the [Available Rules table](#available-rules) details what are the `golint` rules).
When a configuration is provided, only rules in the configuration are enabled.

#### MegaLinter

[MegaLinter](https://megalinter.io/), an open-source linters aggregator for CI, runs `revive` on Go projects out of the box.
See its [revive documentation page](https://megalinter.io/latest/descriptors/go_revive/) for configuration details.

### Command Line Flags

`revive` accepts the following command line parameters:

- `-config [PATH]` - path to the config file in TOML format, defaults to `$HOME/revive.toml` if present.
- `-exclude [PATTERN]` - pattern for files/directories/packages to be excluded for linting.
You can specify the files you want to exclude for linting either as package name (i.e. `github.com/mgechev/revive`),
list them as individual files (i.e. `file.go`), directories (i.e. `./foo/...`), or any combination of the three.
If no exclusion patterns are specified, `vendor/...` will be excluded by default.
- `-formatter [NAME]` - formatter to be used for the output. The currently available formatters are:

  - `default` - will output the failures the same way that `golint` does.
  - `json` - outputs the failures in JSON format.
  - `ndjson` - outputs the failures as a stream in newline delimited JSON (NDJSON) format.
  - `friendly` - outputs the failures when found. Shows the summary of all the failures.
  - `stylish` - formats the failures in a table. Keep in mind that it doesn't stream the output so it might be perceived as slower compared to others.
  - `checkstyle` - outputs the failures in XML format compatible with that of Java's [Checkstyle](https://checkstyle.org/).
- `-max_open_files` -  maximum number of open files at the same time. Defaults to unlimited.
- `-set_exit_status` - set exit status to 1 if any issues are found, overwrites `error-code` and `warning-code` in config.
- `-version` - get revive version.

### Sample Invocations

```shell
revive -config revive.toml -exclude file1.go -exclude file2.go -formatter friendly github.com/mgechev/revive package/...
```

- The command above will use the configuration from `revive.toml`
- `revive` will ignore `file1.go` and `file2.go`
- The output will be formatted with the `friendly` formatter
- The linter will analyze `github.com/mgechev/revive` and the files in `package`

### Comment Directives

Using comments, you can disable the linter for the entire file or only a range of lines:

```go
//revive:disable

func Public() {}

//revive:enable
```

The snippet above, will disable `revive` between the `revive:disable` and `revive:enable` comments.
If you skip `revive:enable`, the linter will be disabled for the rest of the file.

With `revive:disable-next-line` and `revive:disable-line` you can disable `revive` on a particular code line.

You can do the same on a rule level. In case you want to disable only a particular rule, you can use:

```go
//revive:disable:unexported-return
func Public() private {
	return private
}

//revive:enable:unexported-return
```

This way, `revive` will not warn you that you're returning an object of an unexported type, from an exported function.

You can document why you disable the linter by adding a trailing text in the directive, for example

```go
//revive:disable Until the code is stable
```

```go
//revive:disable:cyclomatic High complexity score but easy to understand
```

You can also configure `revive` to enforce documenting linter disabling directives by adding

```toml
[directive.specify-disable-reason]
```

in the configuration. You can set the severity (defaults to _warning_) of the violation of this directive

```toml
[directive.specify-disable-reason]
severity = "error"
```

Similarly, you can enforce specifying rule names in disable directives by adding

```toml
[directive.specify-disable-rule]
```

This forbids "naked" `//revive:disable`, `//revive:disable-line`, and `//revive:disable-next-line` directives
and requires specifying at least one rule name, e.g. `//revive:disable:cyclomatic` or `//revive:disable-next-line:cyclomatic`.
You can set the severity (defaults to _warning_) of the violation of this directive

```toml
[directive.specify-disable-rule]
severity = "error"
```

### Configuration

`revive` can be configured with a TOML file. Here's a sample configuration with an explanation of the individual properties:

```toml
# When set to false, ignores files with "GENERATED" header, similar to golint
ignore-generated-header = true

# Sets the default severity to "warning"
severity = "warning"

# Sets the default failure confidence. This means that linting errors
# with less than 0.8 confidence will be ignored.
confidence = 0.8

# Sets the error code for failures with the "error" severity
error-code = 0

# Sets the error code for failures with severity "warning"
warning-code = 0

# Configuration of the `cyclomatic` rule. Here we specify that
# the rule should fail if it detects code with higher complexity than 10.
[rule.cyclomatic]
arguments = [10]

# Sets the severity of the `package-comments` rule to "error".
[rule.package-comments]
severity = "error"
```

By default `revive` will enable only the linting rules that are named in the configuration file.
For example, the previous configuration file makes `revive` to enable only _cyclomatic_ and _package-comments_ linting rules.

To enable default rules you need to use:

```toml
enable-default-rules = true
```

This will enable all rules available in `golint` and use their default configuration (i.e. the way they are hardcoded in `golint`).
The default configuration of `revive` can be found at `defaults.toml`.

To enable all available rules you need to add:

```toml
enable-all-rules = true
```

This will enable all available rules no matter what rules are named in the configuration file.

Options `enable-all-rules` and `enable-default-rules` cannot be combined.

To disable a rule, you simply mark it as disabled in the configuration.
For example:

```toml
[rule.line-length-limit]
disabled = true
```

When enabling all rules you still need/can provide specific configurations for rules.
The following file is an example configuration where all rules are enabled, except for those that are explicitly disabled,
and some rules are configured with particular arguments:

```toml
severity = "warning"
confidence = 0.8
error-code = 0
warning-code = 0

# Enable all available rules
enable-all-rules = true

# Disabled rules
[rule.blank-imports]
disabled = true
[rule.file-header]
disabled = true
[rule.max-public-structs]
disabled = true
[rule.line-length-limit]
disabled = true
[rule.function-length]
disabled = true
[rule.banned-characters]
disabled = true

# Rule tuning
[rule.argument-limit]
arguments = [5]
[rule.cyclomatic]
arguments = [10]
[rule.cognitive-complexity]
arguments = [7]
[rule.function-result-limit]
arguments = [3]
[rule.error-strings]
arguments = ["mypackage.Error"]
```

### Custom Configuration

```shell
revive -config config.toml -formatter friendly github.com/mgechev/revive
```

This will use `config.toml`, the `friendly` formatter, and will run linting over the `github.com/mgechev/revive` package.

### Recommended Configuration

The following snippet contains the recommended `revive` configuration that you can use in your project:

```toml
ignore-generated-header = false
severity = "warning"
confidence = 0.8
error-code = 0
warning-code = 0

[rule.blank-imports]
[rule.context-as-argument]
[rule.context-keys-type]
[rule.dot-imports]
[rule.error-return]
[rule.error-strings]
[rule.error-naming]
[rule.exported]
[rule.increment-decrement]
[rule.var-naming]
[rule.var-declaration]
[rule.package-comments]
[rule.range]
[rule.receiver-naming]
[rule.time-naming]
[rule.unexported-return]
[rule.indent-error-flow]
[rule.errorf]
[rule.empty-block]
[rule.superfluous-else]
[rule.unused-parameter]
[rule.unreachable-code]
[rule.redefines-builtin-id]
```

### Rule-level file excludes

You also can setup custom excludes for each rule.

It's an alternative for the global `-exclude` program arg.

```toml
ignore-generated-header = false
severity = "warning"
confidence = 0.8
error-code = 0
warning-code = 0

[rule.blank-imports]
exclude = ["**/*.pb.go"]
[rule.context-as-argument]
exclude = ["src/somepkg/*.go", "TEST"]
```

You can use the following exclude patterns

1. full paths to files `src/pkg/mypkg/some.go`
2. globs `src/**/*.pb.go`
3. regexes (should have prefix ~) `~\.(pb|auto|generated)\.go$`
4. well-known `TEST` (same as `**/*_test.go`)
5. special cases:
  a. `*` and `~` patterns exclude all files (same effect as disabling the rule)
  b. `""` (empty) pattern excludes nothing

> NOTE: do not mess with `exclude` that can  be used at the top level of TOML file, that means "exclude package patterns", not "exclude file patterns"

## Available Rules

List of all [available rules](./RULES_DESCRIPTIONS.md).

| Name                  | Config | Description                                                      | Go version | Revive version | `golint` | Typed |
| --------------------- | :----: | :--------------------------------------------------------------- | :--------: | :------------: | :------: | :---: |
| [`add-constant`](./RULES_DESCRIPTIONS.md#add-constant)        |  map   | Suggests using constant for magic numbers and string literals    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`argument-limit`](./RULES_DESCRIPTIONS.md#argument-limit)      |  int (defaults to 8)  | Specifies the maximum number of arguments a function can receive | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`atomic`](./RULES_DESCRIPTIONS.md#atomic)              |  n/a   | Check for common mistaken usages of the `sync/atomic` package    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`banned-characters`](./RULES_DESCRIPTIONS.md#banned-characters)          |  []string (defaults to []string{})   |  Checks banned characters in identifiers | 1.0 | [v1.1.3](https://github.com/revive-lint/revive/releases/tag/v1.1.3) |    no    |  no   |
| [`bare-return`](./RULES_DESCRIPTIONS.md#bare-return) | n/a  | Warns on bare returns   | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`blank-imports`](./RULES_DESCRIPTIONS.md#blank-imports)       |  n/a   | Disallows blank imports                                          | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`bool-literal-in-expr`](./RULES_DESCRIPTIONS.md#bool-literal-in-expr)|  n/a   | Suggests removing Boolean literals from logic expressions        | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`call-to-gc`](./RULES_DESCRIPTIONS.md#call-to-gc)   | n/a    | Warns on explicit call to the garbage collector    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`cognitive-complexity`](./RULES_DESCRIPTIONS.md#cognitive-complexity)          |  int (defaults to 7) | Sets restriction for maximum Cognitive complexity.              | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`comment-spacings`](./RULES_DESCRIPTIONS.md#comment-spacings)          |  []string   |  Warns on malformed comments | 1.0 | [v1.2.5](https://github.com/revive-lint/revive/releases/tag/v1.2.5) |    no    |  no   |
| [`comments-density`](./RULES_DESCRIPTIONS.md#comments-density) |  int (defaults to 0)  | Enforces a minimum comment / code relation | 1.0 | [v1.3.8](https://github.com/revive-lint/revive/releases/tag/v1.3.8) |    no    |  no   |
| [`confusing-naming`](./RULES_DESCRIPTIONS.md#confusing-naming)    |  n/a   | Warns on methods with names that differ only by capitalization   | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`confusing-results`](./RULES_DESCRIPTIONS.md#confusing-results)   |  n/a   | Suggests to name potentially confusing function results          | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`constant-logical-expr`](./RULES_DESCRIPTIONS.md#constant-logical-expr)   |  n/a   | Warns on constant logical expressions                        | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`context-as-argument`](./RULES_DESCRIPTIONS.md#context-as-argument) |  n/a   | `context.Context` should be the first argument of a function.    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`context-keys-type`](./RULES_DESCRIPTIONS.md#context-keys-type)   |  n/a   | Disallows the usage of basic types in `context.WithValue`.       | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  yes  |
| [`cyclomatic`](./RULES_DESCRIPTIONS.md#cyclomatic)          |  int (defaults to 10)   | Sets restriction for maximum Cyclomatic complexity.              | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`datarace`](./RULES_DESCRIPTIONS.md#datarace)          |  n/a   |  Spots potential dataraces | 1.0 | [v1.2.2](https://github.com/revive-lint/revive/releases/tag/v1.2.2) |    no    |  no   |
| [`deep-exit`](./RULES_DESCRIPTIONS.md#deep-exit)           |  n/a   | Looks for program exits in funcs other than `main()` or `init()` | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`defer`](./RULES_DESCRIPTIONS.md#defer)          |  map   |  Warns on some [defer gotchas](https://blog.learngoprogramming.com/5-gotchas-of-defer-in-go-golang-part-iii-36a1ab3d6ef1)       | 1.0 | [v1.0.3](https://github.com/revive-lint/revive/releases/tag/v1.0.3) |    no    |  no   |
| [`dot-imports`](./RULES_DESCRIPTIONS.md#dot-imports)         |  n/a   | Forbids `.` imports.                                             | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`duplicated-imports`](./RULES_DESCRIPTIONS.md#duplicated-imports) | n/a  | Looks for packages that are imported two or more times   | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`early-return`](./RULES_DESCRIPTIONS.md#early-return)          | []string   | Spots if-then-else statements where the predicate may be inverted to reduce nesting | 1.0 | [v1.0.3](https://github.com/revive-lint/revive/releases/tag/v1.0.3) |    no    |  no   |
| [`empty-block`](./RULES_DESCRIPTIONS.md#empty-block)         |  n/a   | Warns on empty code blocks                                       | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`empty-lines`](./RULES_DESCRIPTIONS.md#empty-lines)   | n/a | Warns when there are heading or trailing newlines in a block              | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`epoch-naming`](./RULES_DESCRIPTIONS.md#epoch-naming)      |  n/a   | Enforces naming conventions for epoch time variables       | 1.0 | [v1.14.0](https://github.com/revive-lint/revive/releases/tag/v1.14.0) |    no    |  yes   |
| [`enforce-map-style`](./RULES_DESCRIPTIONS.md#enforce-map-style) |  string (defaults to "any")  |  Enforces consistent usage of `make(map[type]type)` or `map[type]type{}` for map initialization. Does not affect `make(map[type]type, size)` constructions. | 1.0 | [v1.3.4](https://github.com/revive-lint/revive/releases/tag/v1.3.4) |    no    |  no   |
| [`enforce-repeated-arg-type-style`](./RULES_DESCRIPTIONS.md#enforce-repeated-arg-type-style) |  string (defaults to "any")  |  Enforces consistent style for repeated argument and/or return value types. | 1.0 | [v1.3.5](https://github.com/revive-lint/revive/releases/tag/v1.3.5) |    no    |  no   |
| [`enforce-slice-style`](./RULES_DESCRIPTIONS.md#enforce-slice-style) |  string (defaults to "any")  |  Enforces consistent usage of `make([]type, 0)` or `[]type{}` for slice initialization. Does not affect `make(map[type]type, non_zero_len, or_non_zero_cap)` constructions. | 1.0 | [v1.3.5](https://github.com/revive-lint/revive/releases/tag/v1.3.5) |    no    |  no   |
| [`enforce-switch-style`](./RULES_DESCRIPTIONS.md#enforce-switch-style) |  []string (defaults to enforce occurrence and position) |  Enforces consistent usage of `default` on `switch` statements. | 1.0 | [v1.11.0](https://github.com/revive-lint/revive/releases/tag/v1.11.0) |    no    |  no   |
| [`error-naming`](./RULES_DESCRIPTIONS.md#error-naming)        |  n/a   | Naming of error variables.                                       | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`error-return`](./RULES_DESCRIPTIONS.md#error-return)        |  n/a   | The error return parameter should be last.                       | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`error-strings`](./RULES_DESCRIPTIONS.md#error-strings)       |  []string   | Conventions around error strings.                                | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`errorf`](./RULES_DESCRIPTIONS.md#errorf)              |  n/a   | Should replace `errors.New(fmt.Sprintf())` with `fmt.Errorf()`   | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  yes  |
| [`exported`](./RULES_DESCRIPTIONS.md#exported)            |  []string   | Naming and commenting conventions on exported symbols.           | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`file-header`](./RULES_DESCRIPTIONS.md#file-header)         | string (defaults to none)| Header which each file should have.                              | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`file-length-limit`](./RULES_DESCRIPTIONS.md#file-length-limit) | map (optional)| Enforces a maximum number of lines per file | 1.0 | [v1.5.0](https://github.com/revive-lint/revive/releases/tag/v1.5.0) |    no    |  no   |
| [`filename-format`](./RULES_DESCRIPTIONS.md#filename-format) | regular expression (optional) | Enforces the formatting of filenames | 1.0 | [v1.5.0](https://github.com/revive-lint/revive/releases/tag/v1.5.0) |   no    |  no   |
| [`flag-parameter`](./RULES_DESCRIPTIONS.md#flag-parameter)      |  n/a   | Warns on boolean parameters that create a control coupling       | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`forbidden-call-in-wg-go`](./RULES_DESCRIPTIONS.md#forbidden-call-in-wg-go)  |  n/a   | Warns on forbidden calls inside calls to wg.Go | >= 1.25 | [v1.13.0](https://github.com/revive-lint/revive/releases/tag/v1.13.0) |    no    |  no   |
| [`function-length`](./RULES_DESCRIPTIONS.md#function-length)          |  int, int (defaults to 50 statements, 75 lines)   |  Warns on functions exceeding the statements or lines max | 1.0 | [v1.0.5](https://github.com/revive-lint/revive/releases/tag/v1.0.5) |    no    |  no   |
| [`function-result-limit`](./RULES_DESCRIPTIONS.md#function-result-limit) |  int (defaults to 3)| Specifies the maximum number of results a function can return    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`get-return`](./RULES_DESCRIPTIONS.md#get-return)          |  n/a   | Warns on getters that do not yield any result                    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`identical-branches`](./RULES_DESCRIPTIONS.md#identical-branches)          |  n/a   | Spots if-then-else statements with identical `then` and `else` branches       | 1.0 | [v1.0.3](https://github.com/revive-lint/revive/releases/tag/v1.0.3) |    no    |  no   |
| [`identical-ifelseif-branches`](./RULES_DESCRIPTIONS.md#identical-ifelseif-branches)          |  n/a   | Spots `if ... else if` chains with identical branches.       | 1.0 | [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0) |    no    |  no   |
| [`identical-ifelseif-conditions`](./RULES_DESCRIPTIONS.md#identical-ifelseif-conditions)          |  n/a   | Spots identical conditions in  `if ... else if` chains.       | 1.0 | [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0) |    no    |  no   |
| [`identical-switch-branches`](./RULES_DESCRIPTIONS.md#identical-switch-branches)          |  map   | Spots `switch` with identical branches.       | 1.0 | [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0) |    no    |  no   |
| [`identical-switch-conditions`](./RULES_DESCRIPTIONS.md#identical-switch-conditions)          |  n/a   | Spots identical conditions in case clauses of `switch` statements.       | 1.0 | [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0) |    no    |  no   |
| [`if-return`](./RULES_DESCRIPTIONS.md#if-return)           |  n/a   | Redundant if when returning an error.                            | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   no    |  no   |
| [`import-alias-naming`](./RULES_DESCRIPTIONS.md#import-alias-naming)         | string or map[string]string (defaults to allow regex pattern `^[a-z][a-z0-9]{0,}$`) | Conventions around the naming of import aliases.                              | 1.0 | [v1.3.4](https://github.com/revive-lint/revive/releases/tag/v1.3.4) |    no    |  no   |
| [`import-shadowing`](./RULES_DESCRIPTIONS.md#import-shadowing)   | n/a    | Spots identifiers that shadow an import    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`imports-blocklist`](./RULES_DESCRIPTIONS.md#imports-blocklist)   | []string | Disallows importing the specified packages                     | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`increment-decrement`](./RULES_DESCRIPTIONS.md#increment-decrement) |  n/a   | Use `i++` and `i--` instead of `i += 1` and `i -= 1`.            | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`indent-error-flow`](./RULES_DESCRIPTIONS.md#indent-error-flow)   |  []string   | Prevents redundant else statements.                              | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`inefficient-map-lookup`](./RULES_DESCRIPTIONS.md#inefficient-map-lookup)   |  n/a  | Spots iterative searches for a key in a map           | 1.0 | [v1.13.0](https://github.com/revive-lint/revive/releases/tag/v1.13.0) |   no    |  yes   |
| [`line-length-limit`](./RULES_DESCRIPTIONS.md#line-length-limit)   | int or map | Specifies the maximum number of characters in a line, with optional per-line exclusions | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`marshal-receiver`](./RULES_DESCRIPTIONS.md#marshal-receiver) |  n/a   | Checks receiver type consistency for common marshal/unmarshal methods. | 1.0 | [v1.16.0](https://github.com/revive-lint/revive/releases/tag/v1.16.0) |    no    |  no   |
| [`max-control-nesting`](./RULES_DESCRIPTIONS.md#max-control-nesting) |  int (defaults to 5)  | Sets restriction for maximum nesting of control structures. | 1.0 | [v1.3.7](https://github.com/revive-lint/revive/releases/tag/v1.3.7) |    no    |  no   |
| [`max-public-structs`](./RULES_DESCRIPTIONS.md#max-public-structs)  |  int (defaults to 5)  | The maximum number of public structs in a file.                  | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`modifies-parameter`](./RULES_DESCRIPTIONS.md#modifies-parameter)  |  n/a   | Warns on assignments to function parameters                      | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`modifies-value-receiver`](./RULES_DESCRIPTIONS.md#modifies-value-receiver) |  n/a   | Warns on assignments to value-passed method receivers        | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  yes  |
| [`multiline-if-init`](./RULES_DESCRIPTIONS.md#multiline-if-init)    |  n/a   | Flags `if` statements whose init clause spans multiple lines    | 1.0 | [v1.16.0](https://github.com/revive-lint/revive/releases/tag/v1.16.0) |    no    |  no   |
| [`nested-structs`](./RULES_DESCRIPTIONS.md#nested-structs)          |  n/a   |  Warns on structs within structs | 1.0 | [v1.0.8](https://github.com/revive-lint/revive/releases/tag/v1.0.8) |    no    |  no   |
| [`optimize-operands-order`](./RULES_DESCRIPTIONS.md#optimize-operands-order)          |  n/a   |  Checks inefficient conditional expressions | 1.0 | [v1.1.3](https://github.com/revive-lint/revive/releases/tag/v1.1.3) |    no    |  no   |
| [`package-comments`](./RULES_DESCRIPTIONS.md#package-comments)    |  n/a   | Package commenting conventions.                                  | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`package-naming`](./RULES_DESCRIPTIONS.md#package-naming)    |  map   | Checks that package names follow Go conventions and best practices | 1.0 | [v1.15.0](https://github.com/revive-lint/revive/releases/tag/v1.15.0) |   no    |  no   |
| [`package-directory-mismatch`](./RULES_DESCRIPTIONS.md#package-directory-mismatch)    | string | Checks that package name matches containing directory name     | 1.0 | [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0) |   no    |  no   |
| [`range-val-address`](./RULES_DESCRIPTIONS.md#range-val-address)|  n/a   | Warns if address of range value is used dangerously | < 1.22 | [v1.0.2](https://github.com/revive-lint/revive/releases/tag/v1.0.2) |    no    |  yes   |
| [`range-val-in-closure`](./RULES_DESCRIPTIONS.md#range-val-in-closure)|  n/a   | Warns if range value is used in a closure dispatched as goroutine| < 1.22 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`range`](./RULES_DESCRIPTIONS.md#range)               |  n/a   | Prevents redundant variables when iterating over a collection.   | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`receiver-naming`](./RULES_DESCRIPTIONS.md#receiver-naming)     |  map (optional)   | Conventions around the naming of receivers.                      | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`redefines-builtin-id`](./RULES_DESCRIPTIONS.md#redefines-builtin-id)|  n/a   | Warns on redefinitions of builtin identifiers                    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`redundant-build-tag`](./RULES_DESCRIPTIONS.md#redundant-build-tag) | n/a   | Warns about build tags that are redundant: a `// +build` line, a `go1.X` constraint already guaranteed by go.mod, or a constraint implied by the file name | >= 1.17 | [v1.6.0](https://github.com/revive-lint/revive/releases/tag/v1.6.0) |   no    |  no   |
| [`redundant-import-alias`](./RULES_DESCRIPTIONS.md#redundant-import-alias)          |  n/a   |  Warns on import aliases matching the imported package name | 1.0 | [v1.3.3](https://github.com/revive-lint/revive/releases/tag/v1.3.3) |    no    |  no   |
| [`redundant-test-main-exit`](./RULES_DESCRIPTIONS.md#redundant-test-main-exit) |  n/a   | Suggests removing `Exit` call in `TestMain` function for test files | >= 1.15 | [v1.6.1](https://github.com/revive-lint/revive/releases/tag/v1.6.1) |    no    |  no   |
| [`string-format`](./RULES_DESCRIPTIONS.md#string-format)          |  map   | Warns on specific string literals that fail one or more user-configured regular expressions            | 1.0 | [v1.0.7](https://github.com/revive-lint/revive/releases/tag/v1.0.7) |    no    |  no   |
| [`string-of-int`](./RULES_DESCRIPTIONS.md#string-of-int)          |  n/a   | Warns on suspicious casts from int to string            | 1.0 | [v1.0.2](https://github.com/revive-lint/revive/releases/tag/v1.0.2) |    no    |  yes   |
| [`struct-tag`](./RULES_DESCRIPTIONS.md#struct-tag)          |  []string   | Checks common struct tags like `json`, `xml`, `yaml`               | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`superfluous-else`](./RULES_DESCRIPTIONS.md#superfluous-else)    |  []string   | Prevents redundant else statements (extends [`indent-error-flow`](./RULES_DESCRIPTIONS.md#indent-error-flow)) | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`time-date`](./RULES_DESCRIPTIONS.md#time-date)         |  n/a   | Reports bad usage of `time.Date`.                 | 1.0 | [v1.10.0](https://github.com/revive-lint/revive/releases/tag/v1.10.0) |   no    |  no   |
| [`time-equal`](./RULES_DESCRIPTIONS.md#time-equal)         |  n/a   | Suggests to use `time.Time.Equal` instead of `==` and `!=` for equality check time.                 | 1.0 | [v1.1.3](https://github.com/revive-lint/revive/releases/tag/v1.1.3) |   no    |  yes  |
| [`time-naming`](./RULES_DESCRIPTIONS.md#time-naming)         |  n/a   | Conventions around the naming of time variables.                 | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  yes  |
| [`unchecked-type-assertion`](./RULES_DESCRIPTIONS.md#unchecked-type-assertion)         |  n/a   | Disallows type assertions without checking the result.                 | 1.0 | [v1.3.4](https://github.com/revive-lint/revive/releases/tag/v1.3.4) |   no    |  no   |
| [`unconditional-recursion`](./RULES_DESCRIPTIONS.md#unconditional-recursion)          |  n/a   | Warns on function calls that will lead to (direct) infinite recursion | 1.0 | [v1.0.3](https://github.com/revive-lint/revive/releases/tag/v1.0.3) |    no    |  no   |
| [`unexported-naming`](./RULES_DESCRIPTIONS.md#unexported-naming)          |  n/a   |  Warns on wrongly named un-exported symbols       | 1.0 | [v1.0.3](https://github.com/revive-lint/revive/releases/tag/v1.0.3) |    no    |  no   |
| [`unexported-return`](./RULES_DESCRIPTIONS.md#unexported-return)   |  n/a   | Warns when a public return is from unexported type.              | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  yes  |
| [`unhandled-error`](./RULES_DESCRIPTIONS.md#unhandled-error)   | []string   | Warns on unhandled errors returned by function calls    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  yes   |
| [`unnecessary-if`](./RULES_DESCRIPTIONS.md#unnecessary-if)    |  n/a   | Identifies `if-else` statements that can be replaced by simpler statements  | 1.0 | [v1.13.0](https://github.com/revive-lint/revive/releases/tag/v1.13.0) |    no    |  no   |
| [`unnecessary-format`](./RULES_DESCRIPTIONS.md#unnecessary-format)    |  n/a   | Identifies calls to formatting functions where the format string does not contain any formatting verbs          | 1.0 | [v1.10.0](https://github.com/revive-lint/revive/releases/tag/v1.10.0) |    no    |  no   |
| [`unnecessary-stmt`](./RULES_DESCRIPTIONS.md#unnecessary-stmt)    |  n/a   | Suggests removing or simplifying unnecessary statements          | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`unreachable-code`](./RULES_DESCRIPTIONS.md#unreachable-code)    |  n/a   | Warns on unreachable code                                        | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`unsecure-url-scheme`](./RULES_DESCRIPTIONS.md#unsecure-url-scheme)          |  n/a   |  Checks for usage of potentially unsecure URL schemes. | 1.0 | [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0) |    no    |  no   |
| [`unused-parameter`](./RULES_DESCRIPTIONS.md#unused-parameter)    |  n/a   | Suggests to rename or remove unused function parameters          | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`unused-receiver`](./RULES_DESCRIPTIONS.md#unused-receiver)   | n/a    | Suggests to rename or remove unused method receivers    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |
| [`use-any`](./RULES_DESCRIPTIONS.md#use-any)          |  n/a   |  Proposes to replace `interface{}` with its alias `any` | >= 1.18 | [v1.2.0](https://github.com/revive-lint/revive/releases/tag/v1.2.0) |    no    |  no   |
| [`use-errors-new`](./RULES_DESCRIPTIONS.md#use-errors-new) | n/a   | Spots calls to `fmt.Errorf` that can be replaced by `errors.New` | < 1.26 | [v1.6.0](https://github.com/revive-lint/revive/releases/tag/v1.6.0) |   no    |  no   |
| [`use-fmt-print`](./RULES_DESCRIPTIONS.md#use-fmt-print) | n/a   | Proposes to replace calls to built-in `print` and `println` with their equivalents from `fmt`. | 1.0 | [v1.10.0](https://github.com/revive-lint/revive/releases/tag/v1.10.0) |   no    |  no   |
| [`use-slices-concat`](./RULES_DESCRIPTIONS.md#use-slices-concat) | n/a   | Proposes to replace appends to an empty slice with a call to `slices.Concat`. | >= 1.22 | [v1.17.0](https://github.com/revive-lint/revive/releases/tag/v1.17.0) |   no    |  no   |
| [`use-slices-sort`](./RULES_DESCRIPTIONS.md#use-slices-sort) | n/a   | Proposes to replace calls to `sort.Ints`, `sort.Strings` and the like with their equivalents from `slices` package. | >= 1.21 | [v1.14.0](https://github.com/revive-lint/revive/releases/tag/v1.14.0) |   no    |  no   |
| [`use-waitgroup-go`](./RULES_DESCRIPTIONS.md#use-waitgroup-go)          |  n/a   |  Proposes to replace `wg.Add ... go {... wg.Done ...}` idiom with `wg.Go` | >= 1.25 | [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0) |    no    |  no   |
| [`useless-break`](./RULES_DESCRIPTIONS.md#useless-break)          |  n/a   |  Warns on useless `break` statements in case clauses | 1.0 | [v1.1.0](https://github.com/revive-lint/revive/releases/tag/v1.1.0) |    no    |  no   |
| [`useless-fallthrough`](./RULES_DESCRIPTIONS.md#useless-fallthrough)  |  n/a   |  Warns on useless `fallthrough` statements in case clauses | 1.0 | [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0) |    no    |  no   |
| [`var-declaration`](./RULES_DESCRIPTIONS.md#var-declaration)     |  n/a   | Reduces redundancies around variable declaration.                | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  yes  |
| [`var-naming`](./RULES_DESCRIPTIONS.md#var-naming)          |  allowlist & blocklist of initialisms   | Naming rules.                                                    | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |   yes    |  no   |
| [`waitgroup-by-value`](./RULES_DESCRIPTIONS.md#waitgroup-by-value)  |  n/a   | Warns on functions taking sync.WaitGroup as a by-value parameter | 1.0 | [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0) |    no    |  no   |

 The rules ported from `golint` are left unchanged and indicated in the "golint" column.

 The rules in the "Typed" column use the type system, which depends on the `GOROOT` and `GOPATH` environment variables.
 Therefore, if `revive` binary is built with `-trimpath` or run on GitHub Actions, these rules may not work as expected.
 See [#1277](https://github.com/revive-lint/revive/issues/1277) for details.

## Configurable rules

A lot of rules can be configured. See [RULES_DESCRIPTIONS.md](./RULES_DESCRIPTIONS.md) for the list of all configurable rules and their options.

## Available Formatters

This section lists all the available formatters and provides a screenshot for each one.

### Friendly

![Friendly formatter](/assets/formatter-friendly.png)

### Stylish

![Stylish formatter](/assets/formatter-stylish.png)

### Default

The default formatter produces the same output as `golint`.

![Default formatter](/assets/formatter-default.png)

### Plain

The plain formatter produces the same output as the default formatter and appends the URL to the rule description.

![Plain formatter](/assets/formatter-plain.png)

### Unix

The unix formatter produces the same output as the default formatter but surrounds the rules in `[]`.

![Unix formatter](/assets/formatter-unix.png)

### JSON

The `json` formatter produces output in JSON format.

### NDJSON

The `ndjson` formatter produces output in [`Newline Delimited JSON`](https://github.com/ndjson/ndjson-spec) format.

### Checkstyle

The `checkstyle` formatter produces output in a [Checkstyle-like](https://checkstyle.sourceforge.io/) format.

### SARIF

The `sarif`  formatter produces output in SARIF, for _Static Analysis Results Interchange Format_,
a standard JSON-based format for the output of static analysis tools defined and promoted by [OASIS](https://www.oasis-open.org/).

Current supported version of the standard is [SARIF-v2.1.0](https://docs.oasis-open.org/sarif/sarif/v2.1.0/csprd01/sarif-v2.1.0-csprd01.html
).

## Extensibility

The tool can be extended with custom rules or formatters. This section contains additional information on how to implement such.

To extend the linter with a custom rule you can push it to this repository or use `revive` as a library (see below)

To add a custom formatter you'll have to push it to this repository or fork it.
This is due to the limited `-buildmode=plugin` support which [works only on Linux (with known issues)](https://pkg.go.dev/plugin).

### Writing a Custom Rule

See [DEVELOPING.md](./DEVELOPING.md) for instructions on how to write a custom rule.

#### Using `revive` as a library

If a rule is specific to your use case
(i.e. it is not a good candidate to be added to `revive`'s rule set) you can add it to your linter using `revive` as a linting engine.

The following code shows how to use `revive` in your application.
In the example only one rule is added (`myRule`), of course, you can add as many as you need to.
Your rules can be configured programmatically or with the standard `revive` configuration file.
The full rule set of `revive` is also actionable by your application.

```go
package main

import (
	"github.com/mgechev/revive/cli"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/revivelib"
)

func main() {
	cli.RunRevive(revivelib.NewExtraRule(&myRule{}, lint.RuleConfig{}))
}

type myRule struct{}

func (f myRule) Name() string {
	return "myRule"
}

func (f myRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	// ...
}
```

You can still go further and use `revive` without its CLI, as part of your library, or your CLI:

```go
package mylib

import (
	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/revivelib"
)

// Error checking removed for clarity
func LintMyFile(file string) {
	conf, _ := config.GetConfig("../defaults.toml")

	revive, _ := revivelib.New(
		conf, // Configuration file
		true, // Set exit status
		2048, // Max open files

		// Then add as many extra rules as you need
		revivelib.NewExtraRule(&myRule{}, lint.RuleConfig{}),
	)

	failuresChan, err := revive.Lint(
		revivelib.Include(file),
		revivelib.Exclude("./fixtures"),
		// You can use as many revivelib.Include or revivelib.Exclude as required
	)
	if err != nil {
		panic("Shouldn't have failed: " + err.Error())
	}

	// Now let's return the formatted errors
	failures, exitCode, _ := revive.Format("stylish", failuresChan)

	// failures is the string with all formatted lint error messages
	// exit code is 0 if no errors, 1 if errors (unless config options change it)
	// ... do something with them
}

type myRule struct{}

func (f myRule) Name() string {
	return "myRule"
}

func (f myRule) Apply(*lint.File, lint.Arguments) []lint.Failure {
	// ...
}
```

### Custom Formatter

Each formatter needs to implement the following interface:

```go
type Formatter interface {
	Format(<-chan Failure, Config) (string, error)
	Name() string
}
```

The `Format` method accepts a channel of `Failure` instances and the configuration of the enabled rules.
The `Name()` method should return a string different from the names of the already existing rules.
This string is used when specifying the formatter when invoking the `revive` CLI tool.

For a sample formatter, take a look at [this file](/formatter/json.go).

## Speed Comparison

Compared to `golint`, `revive` performs better because it lints the files for each individual rule into a separate goroutine.
Here's a basic performance benchmark on MacBook Pro Early 2013 run on Kubernetes:

### golint

```shellsession
$ time golint kubernetes/... > /dev/null

real    0m54.837s
user    0m57.844s
sys     0m9.146s
```

### revive's speed

```shellsession
# no type checking
$ time revive -config untyped.toml kubernetes/... > /dev/null

real    0m8.471s
user    0m40.721s
sys     0m3.262s
```

Keep in mind that if you use rules that require type checking, the performance may drop to 2x faster than `golint`:

```shellsession
# type checking enabled
$ time revive kubernetes/... > /dev/null

real    0m26.211s
user    2m6.708s
sys     0m17.192s
```

Currently, type-checking is enabled by default. If you want to run the linter without type-checking, remove all typed rules from the configuration file.

## Overriding colorization detection

By default, `revive` determines whether or not to colorize its output based on whether it's connected to a TTY or not.
This works for most use cases, but may not behave as expected if you use `revive` in a pipeline of commands,
where STDOUT is being piped to another command.

To force colorization, add `REVIVE_FORCE_COLOR=1` to the environment you're running in. For example:

```shell
REVIVE_FORCE_COLOR=1 revive -formatter friendly ./... | tee revive.log
```

## Who uses Revive

<!-- markdownlint-disable MD013 -->

- [`tidb`](https://github.com/pingcap/tidb) - TiDB is a distributed HTAP database compatible with the MySQL protocol
- [`grafana`](https://github.com/grafana/grafana) - The tool for beautiful monitoring and metric analytics & dashboards for Graphite, InfluxDB & Prometheus & More
- [`etcd`](https://github.com/etcd-io/etcd) - Distributed reliable key-value store for the most critical data of a distributed system
- [`cadence`](https://github.com/uber/cadence) - Cadence is a distributed, scalable, durable, and highly available orchestration engine by Uber to execute asynchronous long-running business logic in a scalable and resilient way
- [`ferret`](https://github.com/MontFerret/ferret) - Declarative web scraping
- [`gopass`](https://github.com/gopasspw/gopass) - The slightly more awesome standard unix password manager for teams
- [`gitea`](https://github.com/go-gitea/gitea) - Git with a cup of tea, painless self-hosted git service
- [`excelize`](https://github.com/360EntSecGroup-Skylar/excelize) - Go library for reading and writing Microsoft Excel™ (XLSX) files
- [`aurora`](https://github.com/xuri/aurora) - aurora is a web-based Beanstalk queue server console written in Go
- [`soar`](https://github.com/XiaoMi/soar) - SQL Optimizer And Rewriter
- [`pyroscope`](https://github.com/pyroscope-io/pyroscope) - Continuous profiling platform
- [`gorush`](https://github.com/appleboy/gorush) - A push notification server written in Go (Golang).
- [`dry`](https://github.com/moncho/dry) - dry - A Docker manager for the terminal.
- [`go-echarts`](https://github.com/chenjiandongx/go-echarts) - The adorable charts library for Golang
- [`reviewdog`](https://github.com/reviewdog/reviewdog) - Automated code review tool integrated with any code analysis tools regardless of programming language
- [`rudder-server`](https://github.com/rudderlabs/rudder-server) - Privacy and Security focused Segment-alternative, in Golang and React.
- [`sklearn`](https://github.com/pa-m/sklearn) - A partial port of scikit-learn written in Go.
- [`protoc-gen-doc`](https://github.com/pseudomuto/protoc-gen-doc) - Documentation generator plugin for Google Protocol Buffers.
- [`llvm`](https://github.com/llir/llvm) - Library for interacting with LLVM IR in pure Go.
- [`jenkins-library`](https://github.com/SAP/jenkins-library) - Jenkins shared library for Continuous Delivery pipelines by SAP.
- [`pd`](https://github.com/tikv/pd) - Placement driver for TiKV.
- [`shellhub`](https://github.com/shellhub-io/shellhub) - ShellHub enables teams to easily access any Linux device behind firewall and NAT.
- [`lorawan-stack`](https://github.com/TheThingsNetwork/lorawan-stack) - The Things Network Stack for LoRaWAN V3
- [`gin-jwt`](https://github.com/appleboy/gin-jwt) - This is a JWT middleware for Gin framework.
- [`gofight`](https://github.com/appleboy/gofight) - Testing API Handler written in Golang.
- [`Beaver`](https://github.com/Clivern/Beaver) - A Real Time Messaging Server.
- [`ggz`](https://github.com/go-ggz/ggz) - An URL shortener service written in Golang
- [`Codeac.io`](https://www.codeac.io?ref=revive) - Automated code review service integrates with GitHub, Bitbucket and GitLab (even self-hosted) and helps you fight technical debt.
- [`DevLake`](https://github.com/apache/incubator-devlake) - Apache DevLake is an open-source dev data platform to ingest, analyze, and visualize the fragmented data from DevOps tools，which can distill insights to improve engineering productivity.
- [`checker`](https://github.com/cinar/checker) - Checker helps validating user input through rules defined in struct tags or directly through functions.
- [`milvus`](https://github.com/milvus-io/milvus) - A cloud-native vector database, storage for next generation AI applications.
- [`indicator`](https://github.com/cinar/indicator) - Indicator provides various technical analysis indicators, strategies, and a backtesting framework.

<!-- markdownlint-enable MD013 -->

_Open a PR to add your project_.

## Contributors

### Maintainers

| [<img alt="alexandear" src="https://avatars.githubusercontent.com/u/3228886?v=4&s=100" width="100">](https://github.com/alexandear) | [<img alt="chavacava" src="https://avatars.githubusercontent.com/u/25788468?v=4&s=100" width="100">](https://github.com/chavacava) | [<img alt="denisvmedia" src="https://avatars.githubusercontent.com/u/5462781?v=4&s=100" width="100">](https://github.com/denisvmedia) |
|---|---|---|
| [alexandear](https://github.com/alexandear) | [chavacava](https://github.com/chavacava) | [denisvmedia](https://github.com/denisvmedia) |

Revive was created by [mgechev](https://github.com/mgechev).

### All

This project exists thanks to all the people who contribute.

<a href="https://github.com/revive-lint/revive/graphs/contributors">
  <img alt="All Contributors" src="https://contrib.rocks/image?repo=revive-lint/revive&max=500" />
</a>

## License

MIT
