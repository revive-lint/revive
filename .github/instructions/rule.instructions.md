---
description: 'Checklist for adding, modifying, or reviewing a revive rule'
applyTo: 'rule/**/*.go,test/**/*.go,testdata/**,config/config.go,untyped.toml'
---

# Rule Development Instructions

Use this checklist when writing or reviewing a pull request that adds a rule. For a pull request that only modifies an existing rule,
apply the **Interfaces**, **Failures**, **Typed vs untyped**, and **Tests** sections; the naming and documentation items below
describe the layout for new rules, and a few older rules deviate from them (see the exceptions) — do not request renames in unrelated fixes.
[`rule/argument_limit.go`](../../rule/argument_limit.go) and [`test/argument_limit_test.go`](../../test/argument_limit_test.go)
are the canonical examples; [`AGENTS.md`](../../AGENTS.md) §4 points here.

## Identifier

- The rule has a short, meaningful identifier in `kebab-case` (e.g. `argument-limit`).
- The identifier is unique: no rule registered in `allRules` (`config/config.go`) returns the same `Name()`.
  `TestAllRules` in `config/config_test.go` enforces this; a duplicate would otherwise silently shadow the existing rule.
- `Name()` returns exactly that identifier.

## Files and types

- The implementation is a single file `rule/<rule_name>.go`, where `<rule_name>` is the `snake_case` form of the identifier.
- The file declares a struct named after the identifier in `PascalCase` with a `Rule` suffix (e.g. `argument-limit` → `ArgumentsLimitRule`;
  a plural or singular tweak that reads better is fine).
- Identifier, file name, and struct name stay in lockstep for new rules.
  Legacy exceptions exist (e.g. `flag_param.go` → `flag-parameter`, `unused_param.go` → `unused-parameter`,
  `var_declarations.go` → `var-declaration`, and a few structs without the `Rule` suffix such as `FunctionLength` or `NestedStructs`);
  leave them as they are.

## Interfaces

- The struct implements `lint.Rule`: `Name() string` and `Apply(*lint.File, lint.Arguments) []lint.Failure`.
- If the rule takes arguments, it implements `lint.ConfigurableRule` (`Configure(lint.Arguments) error`) and asserts it at compile time
  next to the type declaration, so a signature typo fails the build instead of silently leaving the rule unconfigured:

  ```golang
  var _ lint.ConfigurableRule = (*NewRule)(nil)
  ```

- `Configure` validates arguments and returns an error instead of panicking.
  New options are passed as a single `map[string]any` argument with named keys, not as positional scalars;
  this keeps rule configuration extensible and self-documenting.
- Defaults are constants in the rule file (see `defaultArgumentsLimit`) and are applied in `Configure` when arguments are missing.
  Bundle-level defaults live in `defaults.toml` / `revive.toml`.
- `Configure` resets every configurable field to its default **before** reading the arguments. Rules in `allRules` are singletons and
  `GetLintingRules` may call `Configure` on the same instance more than once, so an option omitted in a later call must not keep
  the value set by an earlier one.

  ```golang
  func (r *SomeRule) Configure(arguments lint.Arguments) error {
      r.max = defaultMax        // reset scalars, maps and slices here,
      r.ignore = nil            // before any early return,
      r.allow = nil             // and for every field, not just the ones below.

      if len(arguments) < 1 {
          return nil // early return must not skip the reset above
      }
      ...
      r.ignore = append(r.ignore, parsed...) // appending to a non-reset slice keeps stale entries
      ...
  }
  ```

  Assigning a field only when its key is present in the argument map has the same effect: the previous value survives.
- `Apply` runs concurrently across files: it never mutates rule state. State is only set in `Configure`.

## Failures

- Each `lint.Failure` sets `Category` to one of the `lint.FailureCategory*` constants in `lint/failure.go`,
  picking the one used by rules with a similar intent. The test harness fails on an empty category.
- Each `lint.Failure` has a `Confidence` that reflects how certain the detection is: use `1` only when a report cannot be a false positive;
  use a lower value (e.g. `0.8`) when the rule relies on heuristics.
  The default `confidence` threshold is `0.8` (`config.DefaultConfidence`), so failures below it are hidden unless the user lowers the threshold;
  the test harness runs with threshold `0` and will not reveal this. Go below `0.8` only when that is intended.
- Failure messages are concise and actionable.

## Typed vs untyped

- A rule is typed if it uses `file.Pkg.TypeCheck()` or anything from `go/types`.
- An untyped rule must be added to `untyped.toml`, keeping the file sorted. A typed rule must not be listed there.

## Registration

- The rule is appended to `allRules` in `config/config.go` so the CLI can discover it.
- `allRulesCount` in `config/config_test.go` is incremented; the enable-all tests fail otherwise.

## Tests

- `test/<rule_name>_test.go` exists and uses the shared `testRule` harness with the standard `testing` package (no assertion libraries).
- Fixtures live under `testdata/<rule_name>.go`, with `_<variant>.go`, `_test.go`, `.gold` files or a sub-directory as needed;
  they cover both reported and non-reported cases, and any configuration option the rule exposes.
- A configurable rule has a `Test<RuleName>ConfigureResetsPreviousOptions` test that reuses one rule instance: it runs `testRule`
  with every option set, then runs it again with no arguments on a fixture that must pass under the defaults.
  See [`test/file_length_limit_test.go`](../../test/file_length_limit_test.go).

## Documentation

- `README.md`: a row is added to the rules table with the correct `Config`, `Go version`, `golint`, and `Typed` columns.
- `RULES_DESCRIPTIONS.md`: a `## <rule-name>` section is added with `_Go version_`, `_Description_`, `_Configuration_`,
  and a TOML configuration example when the rule is configurable, followed by an `### Examples (<rule-name>)` block
  with a snippet that triggers the rule and one that does not.
- New rows and sections go in alphabetical order relative to their neighbours (the existing files are close to, but not strictly, sorted).
- Tables of contents are generated by `markdown-toc`; they are regenerated, not hand-edited.
