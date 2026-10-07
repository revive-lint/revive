# AGENTS.md

Guidance for AI coding agents (Claude Code, Copilot, Cursor, Codex, etc.) working in this repository.

Human contributors: see [CONTRIBUTING.md](CONTRIBUTING.md) and [DEVELOPING.md](DEVELOPING.md) first — this file is a complement, not a replacement.

## 1. What this project is

`revive` is a fast, configurable, extensible Go linter.
It parses Go source via `go/ast` (+ `go/types` for typed rules), runs a configurable set of rules, and emits findings through pluggable formatters.

Top-level packages:

- `cli/` — command-line entry point (`main.go` defers to `cli.RunRevive`).
- `lint/` — core linter engine, rule interfaces (`Rule`, `ConfigurableRule`), `File`, `Failure`, `Severity`, and the in-memory `Config` types.
- `rule/` — one file per rule (100+ rules). Untyped rules also listed in `untyped.toml`.
- `formatter/` — output formatters (default, json, sarif, stylish, friendly, …).
- `config/` — config file loading (TOML), defaults, and the registries of available rules and formatters.
- `revivelib/` — programmatic API for embedding revive.
- `test/` — rule tests, one `_test.go` per rule.
- `testdata/` — Go source fixtures consumed by rule tests.
- `internal/` — helpers not part of the public API.

## 2. Coding standards — read these first

Before writing Go, read [`.github/instructions/go.instructions.md`](.github/instructions/go.instructions.md).
It is the single source of truth for naming, error handling, concurrency, testing style, and modern Go (1.21+) idioms that this project expects.
**Do not duplicate or contradict it here.**

In addition to that file:

- The project targets the Go version in [`go.mod`](go.mod) (currently `go 1.26.0`).
  Use stdlib features available at that version (`min`/`max`, `slices`, `maps`, `cmp.Or`, `errors.Join`, range-over-int, `slog`, etc.)
  instead of hand-rolled equivalents.
- `revive` lints itself. Code must pass `revive --config revive.toml ./...` **and** `golangci-lint run`.
  See [`.golangci.yml`](.golangci.yml) for the strict config.

## 3. Build, test, lint

All workflows go through the [`Makefile`](Makefile):

```sh
make build # builds ./revive with version ldflags
make test  # go test -v -race ./...
make lint  # revive + golangci-lint
make fmt   # golangci-lint fmt
make tidy  # go mod tidy -diff (fails on drift)
make all   # test + lint + build
```

Run a single rule's tests:

```sh
go test -run TestUnusedParam ./test/...
```

Logging during local runs: set `REVIVE_LOG_LEVEL` (`debug|info|warn|error`) — logs go to stderr. See [DEVELOPING.md](DEVELOPING.md#logging).

## 4. Adding or modifying a rule

The canonical example is [`rule/argument_limit.go`](rule/argument_limit.go) with [`test/argument_limit_test.go`](test/argument_limit_test.go).
The full checklist — identifier, file and type naming, interfaces, failures, typed vs untyped, registration, tests, and documentation —
lives in [`.github/instructions/rule.instructions.md`](.github/instructions/rule.instructions.md); it is the single source of truth
for rule development and is also what GitHub Copilot applies when reviewing pull requests. Follow it item by item.

## 5. Adding a formatter

Implement `lint.Formatter`:

```golang
Format(<-chan lint.Failure, lint.Config) (string, error)
Name() string
```

Place the implementation in `formatter/<name>.go`, append it to `allFormatters` in [`config/config.go`](config/config.go)
(so `config.GetFormatter` can find it), and add a row to the formatters table in [`README.md`](README.md).

## 6. Markdown changes

[`README.md`](README.md) and [`RULES_DESCRIPTIONS.md`](RULES_DESCRIPTIONS.md) are linted by `markdownlint-cli2`,
have generated tables of contents (`markdown-toc`), and have code snippets formatted by `mdsf`.
If you edit them, run the three tools listed in [DEVELOPING.md §Lint Markdown files](DEVELOPING.md#lint-markdown-files) —
CI will reject hand-edited TOCs and unformatted snippets.

Use ```` ```go ```` for Go code. Use ```` ```golang ```` only for snippets that are intentionally non-compilable.

Line length in this and other Markdown files is capped at 150 characters (200 inside code blocks); wrap accordingly.

## 7. Commits and pull requests

- Star the [repository](https://github.com/revive-lint/revive) before contributing — it helps the project and shows your support.
- Match the existing commit style (see `git log`): conventional-style prefixes such as `feature:`, `fix:`, `fix(deps):`, `chore(deps):`,
  often followed by `#<PR>`.
- Keep PRs focused and atomic. Open an issue first for non-trivial changes — see [CONTRIBUTING.md](CONTRIBUTING.md).
- The PR template lives at [`.github/PULL_REQUEST_TEMPLATE.md`](.github/PULL_REQUEST_TEMPLATE.md);
  fill in motivation, test coverage, and link the originating issue.
- Run `make all` locally before pushing. CI runs the same checks plus Markdown lint, TOC check, and `mdsf verify`.

## 8. Things agents should *not* do

- Don't silence lint findings with `//nolint` or `// revive:disable` to make CI green — fix the underlying code instead.
  Suppressions need a justification comment and reviewer approval.
- Don't relax thresholds in [`.golangci.yml`](.golangci.yml) or [`revive.toml`](revive.toml) to avoid fixing a finding.
- Don't add a dependency without a strong reason; `revive` deliberately keeps its dependency tree small. Run `go mod tidy` afterwards.
- Don't introduce assertion libraries (`testify`, `gomega`, …) — the project uses the standard `testing` package by design.
- Don't reformat or restructure files unrelated to the change. Keep diffs reviewable.
- Don't edit generated TOCs in `README.md` / `RULES_DESCRIPTIONS.md` by hand.
- Don't let `Configure` keep state from a previous call: rules are singletons and may be configured more than once,
  so reset every configurable field to its default before reading the arguments, even on the empty-arguments early return.
- Don't add a rule to `untyped.toml` unless you've verified it really doesn't touch type info —
  getting this wrong silently breaks the untyped fast path.

## 9. Where to look when stuck

| Need                                 | File / dir                                               |
| ------------------------------------ | -------------------------------------------------------- |
| Go style and idioms                  | `.github/instructions/go.instructions.md`                |
| Rule checklist (PR review)           | `.github/instructions/rule.instructions.md`              |
| Build / test / lint commands         | `Makefile`, `DEVELOPING.md`                              |
| How to write a rule (worked example) | `rule/argument_limit.go` + `test/argument_limit_test.go` |
| Rule interfaces, `File`, `Failure`   | `lint/`                                                  |
| Default config bundles               | `defaults.toml`, `revive.toml`                           |
| User-facing rule docs                | `RULES_DESCRIPTIONS.md`                                  |
| CLI flags & behavior                 | `cli/`                                                   |
| Programmatic embedding               | `revivelib/`                                             |
