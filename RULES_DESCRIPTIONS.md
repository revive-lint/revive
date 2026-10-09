# Description of available rules

List of all available rules.

<!-- toc -->

- [Configuration options format](#configuration-options-format)
- [add-constant](#add-constant)
- [argument-limit](#argument-limit)
- [atomic](#atomic)
- [banned-characters](#banned-characters)
- [bare-return](#bare-return)
- [blank-imports](#blank-imports)
- [bool-literal-in-expr](#bool-literal-in-expr)
- [call-to-gc](#call-to-gc)
- [cognitive-complexity](#cognitive-complexity)
- [comment-spacings](#comment-spacings)
- [comments-density](#comments-density)
- [confusing-naming](#confusing-naming)
- [confusing-results](#confusing-results)
- [constant-logical-expr](#constant-logical-expr)
- [context-as-argument](#context-as-argument)
- [context-keys-type](#context-keys-type)
- [cyclomatic](#cyclomatic)
- [datarace](#datarace)
- [deep-exit](#deep-exit)
- [defer](#defer)
- [dot-imports](#dot-imports)
- [duplicated-imports](#duplicated-imports)
- [early-return](#early-return)
- [empty-block](#empty-block)
- [empty-lines](#empty-lines)
- [epoch-naming](#epoch-naming)
- [enforce-map-style](#enforce-map-style)
- [enforce-repeated-arg-type-style](#enforce-repeated-arg-type-style)
- [enforce-slice-style](#enforce-slice-style)
- [enforce-switch-style](#enforce-switch-style)
- [error-naming](#error-naming)
- [error-return](#error-return)
- [error-strings](#error-strings)
- [errorf](#errorf)
- [exported](#exported)
- [file-header](#file-header)
- [file-length-limit](#file-length-limit)
- [filename-format](#filename-format)
- [flag-parameter](#flag-parameter)
- [forbidden-call-in-wg-go](#forbidden-call-in-wg-go)
- [function-length](#function-length)
- [function-result-limit](#function-result-limit)
- [get-return](#get-return)
- [identical-branches](#identical-branches)
- [identical-ifelseif-branches](#identical-ifelseif-branches)
- [identical-ifelseif-conditions](#identical-ifelseif-conditions)
- [identical-switch-branches](#identical-switch-branches)
- [identical-switch-conditions](#identical-switch-conditions)
- [if-return](#if-return)
- [import-alias-naming](#import-alias-naming)
- [import-shadowing](#import-shadowing)
- [imports-blocklist](#imports-blocklist)
- [increment-decrement](#increment-decrement)
- [indent-error-flow](#indent-error-flow)
- [inefficient-map-lookup](#inefficient-map-lookup)
- [line-length-limit](#line-length-limit)
- [marshal-receiver](#marshal-receiver)
- [max-control-nesting](#max-control-nesting)
- [max-public-structs](#max-public-structs)
- [modifies-parameter](#modifies-parameter)
- [modifies-value-receiver](#modifies-value-receiver)
- [multiline-if-init](#multiline-if-init)
- [nested-structs](#nested-structs)
- [optimize-operands-order](#optimize-operands-order)
- [package-comments](#package-comments)
- [package-naming](#package-naming)
- [package-directory-mismatch](#package-directory-mismatch)
- [range-val-address](#range-val-address)
- [range-val-in-closure](#range-val-in-closure)
- [range](#range)
- [receiver-naming](#receiver-naming)
- [redefines-builtin-id](#redefines-builtin-id)
- [redundant-build-tag](#redundant-build-tag)
- [redundant-import-alias](#redundant-import-alias)
- [redundant-test-main-exit](#redundant-test-main-exit)
- [string-format](#string-format)
- [string-of-int](#string-of-int)
- [struct-tag](#struct-tag)
- [superfluous-else](#superfluous-else)
- [time-date](#time-date)
- [time-equal](#time-equal)
- [time-naming](#time-naming)
- [unchecked-type-assertion](#unchecked-type-assertion)
- [unconditional-recursion](#unconditional-recursion)
- [unexported-naming](#unexported-naming)
- [unexported-return](#unexported-return)
- [unhandled-error](#unhandled-error)
- [unnecessary-if](#unnecessary-if)
- [unnecessary-format](#unnecessary-format)
- [unnecessary-stmt](#unnecessary-stmt)
- [unreachable-code](#unreachable-code)
- [unsecure-url-scheme](#unsecure-url-scheme)
- [unused-parameter](#unused-parameter)
- [unused-receiver](#unused-receiver)
- [use-any](#use-any)
- [use-errors-new](#use-errors-new)
- [use-fmt-print](#use-fmt-print)
- [use-slices-concat](#use-slices-concat)
- [use-slices-sort](#use-slices-sort)
- [use-waitgroup-go](#use-waitgroup-go)
- [useless-break](#useless-break)
- [useless-fallthrough](#useless-fallthrough)
- [var-declaration](#var-declaration)
- [var-naming](#var-naming)
- [waitgroup-by-value](#waitgroup-by-value)

<!-- tocstop -->

## Configuration options format

By convention, configuration options are documented using the `kebab-case` format (e.g., `max-lit-count`, `allow-strs`, `skip-comments`).
For backward compatibility, `camelCase` (e.g., `maxLitCount`, `allowStrs`, `skipComments`)
and `lowercase` (e.g., `maxlitcount`, `allowstrs`, `skipcomments`) formats are still supported but are deprecated.

## add-constant

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Suggests using constant for [magic numbers](https://en.wikipedia.org/wiki/Magic_number_(programming)#Unnamed_numerical_constants)
and string literals.

_Configuration_:

- `max-lit-count`: (string) maximum number of instances of a string literal that are tolerated before a warning is emitted.
- `allow-strs`: (string) comma-separated list of allowed string literals
- `allow-ints`: (string) comma-separated list of allowed integers
- `allow-floats`: (string) comma-separated list of allowed floats
- `ignore-funcs`: (string) comma-separated list of function names regexp patterns to exclude

Configuration example:

```toml
[rule.add-constant]
arguments = [
  { max-lit-count = "3", allow-strs = "\"\"", allow-ints = "0,1,2", allow-floats = "0.0,0.,1.0,1.,2.0,2.", ignore-funcs = "os\\.*,fmt\\.Println,make" },
]
```

## argument-limit

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Warns when a function receives more parameters than the maximum set by the rule's configuration.
Enforcing a maximum number of parameters helps to keep the code readable and maintainable.

_Configuration_: (int) the maximum number of parameters allowed per function. Default: `8`.

Configuration example:

```toml
[rule.argument-limit]
arguments = [4]
```

## atomic

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Check for commonly mistaken usages of the `sync/atomic` package.

_Configuration_: N/A

## banned-characters

_Go version_: 1.0.

_Revive version_: [v1.1.3](https://github.com/revive-lint/revive/releases/tag/v1.1.3).

_Description_: Checks given banned characters in identifiers (func, var, const). Comments are not checked.

_Configuration_: This rule requires a slice of strings, the characters to ban.

Configuration example:

```toml
[rule.banned-characters]
arguments = ["Ω", "Σ", "σ"]
```

## bare-return

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Warns on bare (a.k.a. naked) returns.

### Examples (bare-return)

Before (violation):

```go
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}
```

After (fixed):

```go
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return x, y
}
```

_Configuration_: N/A

## blank-imports

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: Blank import should be only in a main or test package, or have a comment justifying it.

_Configuration_: N/A

## bool-literal-in-expr

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Using Boolean literals (`true`, `false`) in logic expressions may make the code less readable.
This rule suggests removing Boolean literals from logic expressions.

### Examples (bool-literal-in-expr)

Before (violation):

```go
if attachRequired == true {
  // do something
}

if mustReply == false {
  // do something
}
```

After (fixed):

```go
if attachRequired {
  // do something
}

if !mustReply {
  // do something
}
```

_Configuration_: N/A

## call-to-gc

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Explicitly invoking the garbage collector is, except for specific uses in benchmarking, very dubious.

The garbage collector can be configured through environment variables as [described here](https://pkg.go.dev/runtime).

_Configuration_: N/A

## cognitive-complexity

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: [Cognitive complexity](https://www.sonarsource.com/docs/CognitiveComplexity.pdf) is a measure of how hard code is to understand.
While cyclomatic complexity is good to measure "testability" of the code,
cognitive complexity aims to provide a more precise measure of the difficulty of understanding the code.
Enforcing a maximum complexity per function helps to keep code readable and maintainable.

_Configuration_: (int) the maximum function complexity. Default: `7`.

Configuration example:

```toml
[rule.cognitive-complexity]
arguments = [7]
```

## comment-spacings

_Go version_: 1.0.

_Revive version_: [v1.2.5](https://github.com/revive-lint/revive/releases/tag/v1.2.5).

_Description_: Warns on malformed comments.
Spots comments of the form:

```go
//This is a malformed comment: no space between // and the start of the sentence
```

_Configuration_: ([]string) list of exceptions. For example, to accept comments of the form

```go
//mypragma: activate something
//+optional
```

You need to add both `"mypragma:"` and `"+optional"` in the configuration

The following comment prefixes are allowed by default:

- `//#nosec` — [gosec](https://github.com/securego/gosec) security scanner directive
- Go directive comments matching `//[a-z0-9]+:[a-z0-9]` (e.g. `//nolint:linter`, `//go:generate`, `//revive:disable:rule`),
as well as comments starting with "//line ", "//extern ", and "//export "

Configuration example:

```toml
[rule.comment-spacings]
arguments = ["mypragma:", "+optional"]
```

## comments-density

_Go version_: 1.0.

_Revive version_: [v1.3.8](https://github.com/revive-lint/revive/tree/v1.3.8).

_Description_: Spots files not respecting a minimum value for the [_comments lines density_](https://docs.sonarsource.com/sonarqube/latest/user-guide/metric-definitions/)
metric = _comment lines / (lines of code + comment lines) * 100_

_Configuration_: (int) the minimum expected comments lines density. Default: `0`.

Configuration example:

```toml
[rule.comments-density]
arguments = [15]
```

## confusing-naming

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Methods or fields of `struct` that have names different only by capitalization could be confusing.

_Configuration_: N/A

## confusing-results

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Function or methods that return multiple, no named, values of the same type could induce error.

### Examples (confusing-results)

Before (violation):

```go
// getPos yields the geographical position of this tag.
func (t tag) getPos() (float32, float32)
```

After (fixed):

```go
// getPos yields the geographical position of this tag.
func (t tag) getPos() (longitude float32, latitude float32)
```

_Configuration_: N/A

## constant-logical-expr

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: The rule spots logical expressions that evaluate always to the same value.

_Configuration_: N/A

## context-as-argument

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: By [convention](https://go.dev/wiki/CodeReviewComments#contexts), `context.Context` should be the first parameter of a function.
This rule spots function declarations that do not follow the convention.

_Configuration_:

- `allow-types-before`: (string) comma-separated list of types that may be before 'context.Context'

Configuration example:

```toml
[rule.context-as-argument]
arguments = [
  { allow-types-before = "*testing.T,*github.com/user/repo/testing.Harness" },
]
```

## context-keys-type

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

**_Typed_**

_Description_: Basic types should not be used as a key in `context.WithValue`.

_Configuration_: N/A

## cyclomatic

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: [Cyclomatic complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity) is a measure of code complexity.
Enforcing a maximum complexity per function helps to keep code readable and maintainable.

_Configuration_: (int) the maximum function complexity. Default: `10`.

Configuration example:

```toml
[rule.cyclomatic]
arguments = [3]
```

## datarace

_Go version_: 1.0; behavior changes in 1.22.

_Revive version_: [v1.2.2](https://github.com/revive-lint/revive/releases/tag/v1.2.2).

_Description_: This rule spots potential dataraces caused by goroutines capturing (by-reference) particular identifiers of the function from
which goroutines are created.
The rule is able to spot two of such cases: go-routines capturing named return values, and capturing `for-range` values.

_Configuration_: N/A

## deep-exit

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Packages exposing functions that can stop program execution by exiting are hard to reuse.
This rule looks for program exits in functions other than `main()` or `init()`.

_Configuration_: N/A

## defer

_Go version_: 1.0.

_Revive version_: [v1.0.3](https://github.com/revive-lint/revive/releases/tag/v1.0.3).

_Description_: This rule warns on some common mistakes when using `defer` statement. It currently alerts on the following situations:

<!-- markdownlint-disable MD013 -->

| name              | description                                                                                                                                                                                     |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| call-chain        | even if deferring call-chains of the form `foo()()` is valid, it does not help code understanding (only the last call is deferred)                                                             |
| loop              | deferring inside loops can be misleading (deferred functions are not executed at the end of the loop iteration but of the current function) and it could lead to exhausting the execution stack |
| method-call       | deferring a call to a method can lead to subtle bugs if the method does not have a pointer receiver                                                                                             |
| recover           | calling `recover` outside a deferred function has no effect                                                                                                                                     |
| immediate-recover | calling `recover` at the time a defer is registered, rather than as part of the deferred callback.  e.g. `defer recover()` or equivalent.                                                       |
| return            | returning values from a deferred function has no effect                                                                                                                                         |

<!-- markdownlint-enable MD013 -->

These gotchas are [described here](https://blog.learngoprogramming.com/gotchas-of-defer-in-go-1-8d070894cb01).

_Configuration_: By default, all warnings are enabled but it is possible selectively enable them through configuration.
For example to enable only `call-chain` and `loop`:

Configuration example:

```toml
[rule.defer]
arguments = [["call-chain", "loop"]]
```

## dot-imports

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: Importing with `.` makes the programs much harder to understand because it is unclear whether names belong to the current package or
to an imported package.

More [information here](https://go.dev/wiki/CodeReviewComments#import-dot).

_Configuration_:

- `allowed-packages`: (list of strings) list of allowed dot import packages

Configuration example:

```toml
[rule.dot-imports]
arguments = [
  { allowed-packages = [
    "github.com/onsi/ginkgo/v2",
    "github.com/onsi/gomega",
  ] },
]
```

## duplicated-imports

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: It is possible to unintentionally import the same package twice. This rule looks for packages that are imported two or more times.

_Configuration_: N/A

## early-return

_Go version_: 1.0.

_Revive version_: [v1.0.3](https://github.com/revive-lint/revive/releases/tag/v1.0.3).

_Description_: In Go it is idiomatic to minimize nesting statements, a typical example is to avoid if-then-else constructions.
This rule spots constructions like

```golang
if cond {
	// do something
} else {
	// do other thing
	return ...
}
```

where the `if` condition may be inverted in order to reduce nesting:

```golang
if !cond {
	// do other thing
	return ...
}

// do something
```

### Examples (early-return)

Before (violation):

```go
if hasAccess {
	grantResource()
	logAccess()
} else {
	return errNoAccess
}
```

After (fixed):

```go
if !hasAccess {
	return errNoAccess
}

grantResource()
logAccess()
```

_Configuration_: ([]string) rule flags. Available flags are:

- `preserve-scope`: do not suggest refactorings that would increase variable scope
- `allow-jump`: suggest a new jump (`return`, `continue` or `break` statement) if it could unnest multiple statements.
By default, only relocation of _existing_ jumps (i.e. from the `else` clause) are suggested.

Configuration example:

```toml
[rule.early-return]
arguments = ["preserve-scope", "allow-jump"]
```

## empty-block

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Empty blocks make code less readable and could be a symptom of a bug or unfinished refactoring.

_Configuration_: N/A

### Limitations

The `empty-block` rule skips certain empty `for` loops to avoid false positives, but this may introduce false negatives.

The rule skips the following patterns without flagging them:

```go
// for loop whose only clause is a function-call condition (no Init or Post)
for process() {
    // Intentionally empty - process() does the work
}

// bare for-range loop with no key or value variables
for range ch {
    // Intentionally empty - draining the channel
}
```

Note that **all** bare `for range` loops (with no key or value variables) are skipped, not only channel-draining
ones. This means the rule produces **false negatives** for non-channel iterables:

```go
for range s { // not flagged even though the body is empty (false negative)
}
```

A `for range` that assigns the blank identifier (e.g., `for _ = range s {}`) is **not** bare and **will** be flagged.

However, the rule will produce **false positives** for intentionally empty `for` loops that have `Init` or `Post`
clauses, such as:

```go
// False positive: rule will warn even though this is intentional
for _, c := step(); c; _, c = step() {
    // Loop body is intentionally empty; step() does the work
}

// False positive: rule will warn even though this is intentional
for p := 0; bar(p); p++ {
    // Loop body is intentionally empty; bar(p) does the work
}
```

**Workaround**: If you have intentionally empty `for` loops that the rule flags, you can disable the rule
in-place using a directive comment:

```go
//revive:disable:empty-block
for _, c := step(); c; _, c = step() {
    // Intentionally empty - step() does the work
}
//revive:enable:empty-block
```

The reason for this limitation is that properly detecting whether a `for` loop is intentionally empty requires understanding the
**semantics** of the called functions (whether they have side effects, modify state, etc.), which is beyond the scope of static
analysis that this rule performs.

For more details, see:

- <https://github.com/revive-lint/revive/issues/1622>
- <https://github.com/revive-lint/revive/issues/386>

## empty-lines

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Sometimes `gofmt` is not enough to enforce a common formatting of a code-base;
this rule warns when there are heading or trailing newlines in code blocks.

_Configuration_: N/A

## epoch-naming

_Go version_: 1.0.

_Revive version_: [v1.14.0](https://github.com/revive-lint/revive/releases/tag/v1.14.0).

**_Typed_**

_Description_: Variables initialized with epoch time methods (`time.Now().Unix()`, `time.Now().UnixMilli()`,
`time.Now().UnixMicro()`, `time.Now().UnixNano()`) should have names that clearly indicate their time unit to
prevent confusion and potential bugs when working with different time scales.

This rule enforces that variable names contain appropriate suffixes based on the method used:

- `Unix()`: variable name should end with "Sec", "Second" or "Seconds"
- `UnixMilli()`: variable name should end with "Milli" or "Ms"
- `UnixMicro()`: variable name should end with "Micro", "Microsecond", "Microseconds" or "Us"
- `UnixNano()`: variable name should end with "Nano" or "Ns"

The rule checks variable declarations, short variable declarations (`:=`), and regular assignments (`=`).
The suffix matching is case-insensitive and must appear at the end of the variable name.

### Examples (epoch-naming)

Before (violation):

```go
timestamp := time.Now().Unix()           // unclear which unit
createdAt := time.Now().UnixMilli()      // missing unit indicator
t := time.Now().UnixNano()               // lacks required suffix
```

After (fixed):

```go
timestampSec := time.Now().Unix()        // clearly seconds
createdAtMs := time.Now().UnixMilli()    // clearly milliseconds
tNano := time.Now().UnixNano()           // clearly nanoseconds

// Alternative valid names
createdSeconds := time.Now().Unix()      // full word is fine
updatedMicro := time.Now().UnixMicro()   // microseconds
```

_Configuration_: N/A

## enforce-map-style

_Go version_: 1.0.

_Revive version_: [v1.3.4](https://github.com/revive-lint/revive/releases/tag/v1.3.4).

_Description_: This rule enforces consistent usage of `make(map[type]type)` or `map[type]type{}` for map initialization.
It does not affect `make(map[type]type, size)` constructions as well as `map[type]type{k1: v1}`.

_Configuration_: (string) Specifies the enforced style for map initialization. The options are:

- "any": No enforcement (default).
- "make": Enforces the usage of `make(map[type]type)`.
- "literal": Enforces the usage of `map[type]type{}`.

Configuration example:

```toml
[rule.enforce-map-style]
arguments = ["make"]
```

## enforce-repeated-arg-type-style

_Go version_: 1.0.

_Revive version_: [v1.3.5](https://github.com/revive-lint/revive/releases/tag/v1.3.5).

_Description_: This rule is designed to maintain consistency in the declaration of repeated argument and return value types in Go functions.
It supports three styles: 'any', 'short', and 'full'.
The 'any' style is lenient and allows any form of type declaration.
The 'short' style encourages omitting repeated types for conciseness,
whereas the 'full' style mandates explicitly stating the type for each argument
and return value, even if they are repeated, promoting clarity.

_Configuration (1)_: (string) as a single string, it configures both argument
and return value styles. Accepts 'any', 'short', or 'full' (default: 'any').

_Configuration (2)_: (map[string]any) as a map, allows separate configuration
for function arguments and return values. Valid keys are `func-arg-style` and
`func-ret-val-style`, each accepting 'any', 'short', or 'full'. If a key is not
specified, the default value of 'any' is used.

_Note_: The rule applies checks based on the specified styles. For 'full' style,
it flags instances where types are omitted in repeated arguments or return values.
For 'short' style, it highlights opportunities to omit repeated types for brevity.
Incorrect or unknown configuration values will result in an error.

Example (1):

```toml
[rule.enforce-repeated-arg-type-style]
arguments = ["short"]
```

Example (2):

```toml
[rule.enforce-repeated-arg-type-style]
arguments = [{ func-arg-style = "full", func-ret-val-style = "short" }]
```

## enforce-slice-style

_Go version_: 1.0.

_Revive version_: [v1.3.5](https://github.com/revive-lint/revive/releases/tag/v1.3.5).

_Description_: This rule enforces consistent usage of `make([]type, 0)`, `[]type{}`, `var args []type`, or nil slices for slice initialization.
It does not affect `make([]type, non_zero_len, or_non_zero_cap)` constructions as well as `[]type{v1}`.
Nil slices are always permitted.

_Configuration_: (string) Specifies the enforced style for slice initialization. The options are:

- "any": No enforcement (default).
- "make": Enforces the usage of `make([]type, 0)`.
- "literal": Enforces the usage of `[]type{}`.
- "nil": Enforces nil slices, for example `var s []type` in declarations or `[]type(nil)` in expressions.

Configuration example:

```toml
[rule.enforce-slice-style]
arguments = ["make"]
```

### Examples (enforce-slice-style)

With the `"nil"` style, the suggested fix depends on where the slice appears.
In a variable declaration, a nil slice declaration is suggested:

Before (violation):

```go
var args []string = []string{}
```

After (fixed):

```go
var args []string
```

In any other context (function argument, return value, ...), a nil slice expression is suggested:

Before (violation):

```go
func f() []string {
	return []string{}
}
```

After (fixed):

```go
func f() []string {
	return []string(nil)
}
```

## enforce-switch-style

_Go version_: 1.0.

_Revive version_: [v1.11.0](https://github.com/revive-lint/revive/releases/tag/v1.11.0).

_Description_: This rule enforces consistent usage of `default` on `switch` statements.
It can check for `default` case clause occurrence and/or position in the list of case clauses.

_Configuration_: ([]string) Specifies what to enforced: occurrence and/or position. The, non-mutually exclusive, options are:

- "allow-no-default": allows `switch` without `default` case clause.
- "allow-default-not-last": allows `default` case clause to be not the last clause of the `switch`.

Configuration examples:

To enforce that all `switch` statements have a `default` clause as its the last case clause:

```toml
[rule.enforce-switch-style]
```

To enforce that all `switch` statements have a `default` clause but its position is unimportant:

```toml
[rule.enforce-switch-style]
arguments = ["allow-default-not-last"]
```

To enforce that in all `switch` statements with a `default` clause, the `default` is the last case clause:

```toml
[rule.enforce-switch-style]
arguments = ["allow-no-default"]
```

Notice that a configuration including both options will effectively deactivate the whole rule.

## error-naming

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: By convention, for the sake of readability, variables of type `error` must be named with the prefix `err`.
Unexported error variables should start with `err`, exported ones with `Err`.

### Examples (error-naming)

Before (violation):

```go
import "errors"

var invalidInput = errors.New("invalid input")

var TimeoutError = errors.New("connection timed out")
```

After (fixed):

```go
import "errors"

var errInvalidInput = errors.New("invalid input")

var ErrTimeout = errors.New("connection timed out")
```

_Configuration_: N/A

## error-return

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: By convention, for the sake of readability, the errors should be last in the list of returned values by a function.

### Examples (error-return)

Before (violation):

```go
func readConfig(path string) (error, *Config) {
	// ...
}
```

After (fixed):

```go
func readConfig(path string) (*Config, error) {
	// ...
}
```

_Configuration_: N/A

## error-strings

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: By convention, for better readability, error messages should not be capitalized or end with punctuation or a newline.
By default, the rule analyzes functions for creating errors from `fmt`, `errors`, and `github.com/pkg/errors`.
Optionally, the rule can be configured to analyze user functions that create errors.

More [information here](https://go.dev/wiki/CodeReviewComments#error-strings).

_Configuration_: ([]string) the list of additional error functions to check.
The format of values is `package.FunctionName`.

Configuration example:

```toml
[rule.error-strings]
arguments = ["xerrors.Errorf"]
```

## errorf

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

**_Typed_**

_Description_: It is possible to get a simpler program by replacing `errors.New(fmt.Sprintf())` with `fmt.Errorf()`.
This rule spots that kind of simplification opportunities.

_Configuration_: N/A

## exported

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: Exported function and methods should have comments. This warns on undocumented exported functions and methods.

More [information here](https://go.dev/wiki/CodeReviewComments#doc-comments).

_Configuration_: ([]string) rule flags.
Please notice that without configuration, the default behavior of the rule is that of its `golint` counterpart.
Available flags are:

- `check-private-receivers` enables checking public methods of private types
- `disable-stuttering-check` disables checking for method names that stutter with the package name
  (i.e. avoid failure messages of the form _type name will be used as x.XY by other packages, and that stutters; consider calling this Y_)
- `say-repetitive-instead-of-stutters` replaces the use of the term _stutters_ by _repetitive_ in failure messages
- `check-public-interface` enables checking public method definitions in public interface types
- `disable-checks-on-constants` disables all checks on constant declarations
- `disable-checks-on-functions` disables all checks on function declarations
- `disable-checks-on-methods` disables all checks on method declarations
- `disable-checks-on-types` disables all checks on type declarations
- `disable-checks-on-variables` disables all checks on variable declarations

Configuration example:

```toml
[rule.exported]
arguments = [
  "check-private-receivers",
  "disable-stuttering-check",
  "check-public-interface",
  "disable-checks-on-functions",
]
```

## file-header

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: This rule helps to enforce a common header for all source files in a project by spotting those files that do not have the specified header.

_Configuration_: (string) the header to look for in source files.

Configuration example:

```toml
[rule.file-header]
arguments = ["This is the text that must appear at the top of source files."]
```

## file-length-limit

_Go version_: 1.0.

_Revive version_: [v1.5.0](https://github.com/revive-lint/revive/releases/tag/v1.5.0).

_Description_: This rule enforces a maximum number of lines per file, in order to aid in maintainability and reduce complexity.

_Configuration_:

- `max`: (int) a maximum number of lines in a file. Must be non-negative integers. 0 means the rule is disabled (default `0`);
- `skip-comments`: (bool) if true ignore and do not count lines containing just comments (default `false`);
- `skip-blank-lines`: (bool) if true ignore and do not count lines made up purely of whitespace (default `false`).

Configuration example:

```toml
[rule.file-length-limit]
arguments = [{ max = 100, skip-comments = true, skip-blank-lines = true }]
```

## filename-format

_Go version_: 1.0.

_Revive version_: [v1.5.0](https://github.com/revive-lint/revive/releases/tag/v1.5.0).

_Description_: Enforces conventions on source file names. By default, the rule enforces filenames of the form `^[_A-Za-z0-9][_A-Za-z0-9-]*\.go$`.
Optionally, the rule can be configured to enforce other forms.

_Configuration_: (string) regular expression for source filenames.

Configuration example:

```toml
[rule.filename-format]
arguments = ["^[_a-z][_a-z0-9]*\\.go$"]
```

## flag-parameter

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: If a function controls the flow of another by passing it information on what to do, both functions are said to be [control-coupled](https://en.wikipedia.org/wiki/Coupling_(computer_programming)#Procedural_programming).
Coupling among functions must be minimized for better maintainability of the code.
This rule warns on boolean parameters that create a control coupling.

_Configuration_: N/A

## forbidden-call-in-wg-go

_Go version_: >= 1.25.

_Revive version_: [v1.13.0](https://github.com/revive-lint/revive/releases/tag/v1.13.0).

_Description_: Since Go 1.25, it is possible to create goroutines with the method `WaitGroup.Go`.
The `Go` method calls a function in a new goroutine and adds (`Add`) that task to the WaitGroup.
When the function returns, the task is removed (`Done`) from the WaitGroup.

This rule ensures that functions don't panic as is specified
in the [documentation of `WaitGroup.Go`](https://pkg.go.dev/sync#WaitGroup.Go).

The rule also warns against a common mistake when refactoring legacy code:
accidentally leaving behind a call to `WaitGroup.Done`, which can cause subtle bugs or panics.

### Examples (forbidden-call-in-wg-go)

Legacy code with a call to `wg.Done`:

```go
wg := sync.WaitGroup{}

wg.Add(1)
go func() {
  doSomething()
  wg.Done()
}()

wg.Wait
```

Refactored, incorrect, code:

```go
wg := sync.WaitGroup{}

wg.Go(func() {
  doSomething()
  wg.Done()
})

wg.Wait
```

Fixed code:

```go
wg := sync.WaitGroup{}

wg.Go(func() {
  doSomething()
})

wg.Wait
```

_Configuration_: N/A

## function-length

_Go version_: 1.0.

_Revive version_: [v1.0.5](https://github.com/revive-lint/revive/releases/tag/v1.0.5).

_Description_: Functions too long (with many statements and/or lines) can be hard to understand.

_Configuration_: (int, int) the maximum allowed statements and lines.
Set a value to `0` to disable that specific check; if both values are `0`, the rule is disabled. Default: `50`, `75`.

Configuration example:

```toml
[rule.function-length]
arguments = [10, 0]
```

Will check for functions exceeding 10 statements and will not check the number of lines of functions

## function-result-limit

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Specifies the maximum number of results a function can return.
Functions returning too many results can be hard to understand/use.

_Configuration_: (int) the maximum allowed return values. Default: `3`.

Configuration example:

```toml
[rule.function-result-limit]
arguments = [3]
```

## get-return

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Warns on getters that do not yield any result.
Typically, functions with names prefixed with _Get_ are supposed to return a value.

_Configuration_: N/A

## identical-branches

_Go version_: 1.0.

_Revive version_: [v1.0.3](https://github.com/revive-lint/revive/releases/tag/v1.0.3).

_Description_: An `if-then-else` conditional with identical implementations in both branches is an error.

_Configuration_: N/A

## identical-ifelseif-branches

_Go version_: 1.0.

_Revive version_: [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0).

_Description_: An `if ... else if` chain with identical branches makes maintenance harder
and might be a source of bugs. Duplicated branches should be consolidated in one.

_Configuration_: N/A

## identical-ifelseif-conditions

_Go version_: 1.0.

_Revive version_: [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0).

_Description_: An `if ... else if` chain with identical conditions can lead to
unreachable code and is a potential source of bugs while making the code harder to read and maintain.

_Configuration_: N/A

## identical-switch-branches

_Go version_: 1.0.

_Revive version_: [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0).

_Description_: A `switch` with identical branches makes maintenance harder and might be a source of bugs.
Duplicated branches should be consolidated in one case clause.

_Configuration_: (optional) single map of options (`map[string]any`), provided as the single configuration argument in the rule's arguments array.

- `allow-identical-default`: (bool) If `true`, do not report a `default` case clause that is identical to another case clause. Default: `false`.

### Examples (identical-switch-branches)

When switching over values from an external package, spelling out a fallback that repeats a listed case can be deliberate:
it documents which values are explicitly handled and that everything else lands on the same code.

```go
switch http.SameSite(s) {
case http.SameSiteNoneMode:
	return SameSiteNoneMode
case http.SameSiteLaxMode:
	return SameSiteLaxMode
case http.SameSiteStrictMode:
	return SameSiteStrictMode
default:
	return SameSiteLaxMode
}
```

To keep the rule for identical `case` clauses while allowing that:

```toml
[rule.identical-switch-branches]
arguments = [{ allow-identical-default = true }]
```

## identical-switch-conditions

_Go version_: 1.0.

_Revive version_: [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0).

_Description_: A `switch` statement with cases with the same condition can lead to
unreachable code and is a potential source of bugs while making the code harder to read and maintain.

_Configuration_: N/A

## if-return

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Checking if an error is _nil_ to just after return the error or nil is redundant.

### Examples (if-return)

Before (violation):

```go
func do() error {
	if err := validate(); err != nil {
		return err
	}
	return nil
}
```

After (fixed):

```go
func do() error {
	return validate()
}
```

_Configuration_: N/A

## import-alias-naming

_Go version_: 1.0.

_Revive version_: [v1.3.4](https://github.com/revive-lint/revive/releases/tag/v1.3.4).

_Description_: Aligns with Go's naming conventions, as outlined in the official
[blog post](https://go.dev/blog/package-names). It enforces clear and lowercase import alias names, echoing
the principles of good package naming. Users can follow these guidelines by default or define a custom regex rule.
Importantly, aliases with underscores ("_") are always allowed.

_Configuration_ (1): (`string`) as plain string accepts allow regexp pattern for aliases (default: `^[a-z][a-z0-9]{0,}$`).

_Configuration_ (2): (`map[string]string`) as a map accepts two values:

- for a key `allow-regex` accepts allow regexp pattern
- for a key `deny-regex` deny regexp pattern

_Note_: If both `allow-regex` and `deny-regex` are provided, the alias must comply with both of them.
If none are given (i.e. an empty map), the default value `^[a-z][a-z0-9]{0,}$` for `allow-regex` is used.
Unknown keys will result in an error.

Configuration example (1):

```toml
[rule.import-alias-naming]
arguments = ["^[a-z][a-z0-9]{0,}$"]
```

Configuration example (2):

```toml
[rule.import-alias-naming]
arguments = [{ allow-regex = "^[a-z][a-z0-9]{0,}$", deny-regex = '^v\d+$' }]
```

## import-shadowing

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: In Go it is possible to declare identifiers (packages, structs,
interfaces, parameters, receivers, variables, constants...) that conflict with the
name of an imported package. This rule spots identifiers that shadow an import.

The rule ignores versioned import paths such as `k8s.io/api/core/v1` when `v1` is the package name,
which allows identifiers like `v1`. This is a deliberate trade-off to keep the rule simple.

_Configuration_: N/A

## imports-blocklist

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Warns when importing block-listed packages.

_Configuration_: block-list of package names (or regular expression package names).

Configuration example:

```toml
[rule.imports-blocklist]
arguments = ["crypto/md5", "crypto/sha1", "crypto/**/pkix"]
```

## increment-decrement

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: By convention, for better readability, incrementing an integer variable by 1 is recommended to be done using the `++` operator.
This rule spots expressions like `i += 1` and `i -= 1` and proposes to change them into `i++` and `i--`.

### Examples (increment-decrement)

Before (violation):

```go
i += 1
count -= 1
```

After (fixed):

```go
i++
count--
```

_Configuration_: N/A

## indent-error-flow

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: To improve the readability of code, it is recommended to reduce the indentation as much as possible.
This rule highlights redundant _else-blocks_ that can be eliminated from the code.

More [information here](https://go.dev/wiki/CodeReviewComments#indent-error-flow).

### Examples (indent-error-flow)

Before (violation):

```go
import "log"

if err != nil {
	return err
} else {
	log.Println("no error")
}
```

After (fixed):

```go
import "log"

if err != nil {
	return err
}

log.Println("no error")
```

_Configuration_: ([]string) rule flags. Available flags are:

- `preserve-scope`: do not suggest refactorings that would increase variable scope

Configuration example:

```toml
[rule.indent-error-flow]
arguments = ["preserve-scope"]
```

## inefficient-map-lookup

_Go version_: 1.0.

_Revive version_: [v1.13.0](https://github.com/revive-lint/revive/releases/tag/v1.13.0).

**_Typed_**

_Description_: This rule identifies code that iteratively searches for a key in a map.

This inefficiency is usually introduced when refactoring code from using a slice to a map.
For example if during refactoring the `elements` slice is transformed into a map.

```diff
-       elements             []string
+       elements             map[string]float64
```

and then a loop over `elements` is changed in an obvious but inefficient way:

```diff
-       for _, e := range elements {
+       for e := range elements {
                if e == someStaticValue {
                        // do something
                }
        }
```

Configuration example:

```go
aMap := map[string]bool{}{}
aValue := false

// Inefficient map lookup
for k := range aMap {
	if k == aValue {
		// do something
	}
}

// Simpler and more efficient version
if _, ok := aMap[aValue]; ok {
	// do something
}
```

_Configuration_: N/A

## line-length-limit

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Warns in the presence of code lines longer than a configured maximum.

_Configuration_: the rule accepts either an integer (the maximum line length) or a map of options:

- `max`: (int) maximum line length in characters. Default: `80`.
- `excludes`: ([]string) list of regular expressions; a line matching any of them is not checked. Default: none.

The integer form is a shorthand for setting only `max`.

Configuration examples:

```toml
[rule.line-length-limit]
arguments = [80]
```

```toml
[rule.line-length-limit]
arguments = [{ max = 80, excludes = ['^\s*//go:generate ', 'https?://'] }]
```

### Examples (line-length-limit)

Given the configuration:

```toml
[rule.line-length-limit]
arguments = [{ max = 80, excludes = ['^\s*//go:generate ', 'https?://'] }]
```

for the file:

```go
package example

//go:generate stringer -type=Color -output=color_string.go -linecomment -trimprefix=Color

// Reference: https://example.com/some/really/long/documentation/url/kept/for/context/here

const greeting = "a long enough line of code to comfortably exceed the eighty character limit"
```

only the last line is reported.
The `//go:generate` directive and the line containing a URL match the `excludes` patterns and are skipped, even though both exceed 80 characters.

## marshal-receiver

_Go version_: 1.0.

_Revive version_: [v1.16.0](https://github.com/revive-lint/revive/releases/tag/v1.16.0).

_Description_: Checks receiver type consistency for common marshal/unmarshal methods.
The rule inspects only methods whose names exactly match: `MarshalJSON`, `MarshalText`, `MarshalYAML`, `UnmarshalJSON`, `UnmarshalText`, and `UnmarshalYAML`.
For these methods, it enforces receiver kind only:

- `Marshal*` methods should use a value receiver, and are reported when declared with a pointer receiver.
- `Unmarshal*` methods should use a pointer receiver, and are reported when declared with a value receiver.

This is a name-based, syntactic check.
It does not validate method signatures (parameters or return values) and does not verify whether a method satisfies a specific marshaling interface.

### Examples (marshal-receiver)

Before (violation):

```go
import "encoding/json"

type Person struct {
	Name string
	Age  int
}

func (p *Person) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"name": p.Name,
		"age":  p.Age,
	})
}

func (p Person) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Name = raw["name"].(string)
	p.Age = int(raw["age"].(float64))
	return nil
}
```

After (fixed):

```go
import "encoding/json"

type Person struct {
	Name string
	Age  int
}

// Value receiver — safe, works whether you have Person or *Person.
func (p Person) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"name": p.Name,
		"age":  p.Age,
	})
}

// Pointer receiver — required, must mutate p.
func (p *Person) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Name = raw["name"].(string)
	p.Age = int(raw["age"].(float64))
	return nil
}
```

_Configuration_: N/A

## max-control-nesting

_Go version_: 1.0.

_Revive version_: [v1.3.7](https://github.com/revive-lint/revive/releases/tag/v1.3.7).

_Description_: Warns if nesting level of control structures (`if-then-else`, `for`, `switch`) exceeds a given maximum.

_Configuration_: (int) maximum accepted nesting level of control structures. Default: `5`.

Configuration example:

```toml
[rule.max-control-nesting]
arguments = [3]
```

## max-public-structs

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Packages declaring too many public structs can be hard to understand/use,
and could be a symptom of bad design.

This rule warns on files declaring more than a configured, maximum number of public structs.

_Configuration_: (int) the maximum allowed public structs. Default: `5`.

Configuration example:

```toml
[rule.max-public-structs]
arguments = [3]
```

## modifies-parameter

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: A function that modifies its parameters can be hard to understand.
It can also be misleading if the arguments are passed by value by the caller.
This rule warns when a function modifies one or more of its parameters or when
parameters are passed to functions that modify them (e.g. `slices.Delete`).

_Configuration_: N/A

## modifies-value-receiver

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Typed_**

_Description_: A method that modifies its receiver value can have undesired behavior.
The modification can be also the root of a bug because the actual value receiver could be a copy of that used at the calling site.
This rule warns when a method modifies its receiver.

_Configuration_: N/A

## multiline-if-init

_Go version_: 1.0.

_Revive version_: [v1.16.0](https://github.com/revive-lint/revive/releases/tag/v1.16.0).

_Description_: Flags `if` statements whose init clause spans multiple lines.
The if-init idiom exists for tight one-liners.
When the init wraps across lines, the reader has to visually parse a struct literal or
call chain to find where the initialization ends and the condition begins.
Extract the initialization to a separate statement instead.

### Examples (multiline-if-init)

Before (violation):

```go
if r, err := rec(
	ctx,
	mgr.GetClient(),
	mgr.GetFieldIndexer(),
	mgr.GetEventRecorderFor(fmt.Sprintf("%s-%s-controller", name, options.ManagerName)),
	opts...,
); err != nil {
	return err
}
```

After (fixed):

```go
r, err := rec(
	ctx,
	mgr.GetClient(),
	mgr.GetFieldIndexer(),
	mgr.GetEventRecorderFor(fmt.Sprintf("%s-%s-controller", name, options.ManagerName)),
	opts...,
)
if err != nil {
	return err
}
```

_Configuration_: N/A

## nested-structs

_Go version_: 1.0.

_Revive version_: [v1.0.8](https://github.com/revive-lint/revive/releases/tag/v1.0.8).

_Description_: Packages declaring structs that contain other inline struct definitions can be hard to understand/read for other developers.

_Configuration_: N/A

## optimize-operands-order

_Go version_: 1.0.

_Revive version_: [v1.1.3](https://github.com/revive-lint/revive/releases/tag/v1.1.3).

_Description_: Conditional expressions can be written to take advantage of short circuit evaluation and speed up its average evaluation time
by forcing the evaluation of less time-consuming terms before more costly ones.
This rule spots logical expressions where the order of evaluation of terms seems non optimal.
Please notice that confidence of this rule is low and is up to the user to decide if the suggested rewrite of the expression
keeps the semantics of the original one.

_Configuration_: N/A

Example:

```golang
if isGenerated(content) && !config.IgnoreGeneratedHeader {
```

Swap left and right side :

```golang
if !config.IgnoreGeneratedHeader && isGenerated(content) {
```

## package-comments

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: Packages should have comments. This rule warns on undocumented packages and when packages comments are detached to the `package` keyword.

More [information here](https://go.dev/wiki/CodeReviewComments#package-comments).

_Configuration_: N/A

## package-naming

_Go version_: 1.0.

_Revive version_: [v1.15.0](https://github.com/revive-lint/revive/releases/tag/v1.15.0).

_Description_: This rule checks that package names follow [Go conventions](https://go.dev/blog/package-names) and best practices.
It helps prevent using bad package names and enforces consistent naming patterns.
This rule arose from package naming checks in `var-naming`.

By default, it checks for:

- Package name conventions (no underscores except for test packages, no MixedCaps).
- Bad package names from the official Go blog (e.g., `common`, `util`, `utils`, `misc`, `interfaces`, `types`).
- Package names that conflict with common Go standard library packages (e.g., `http`, `json`, `fmt`).

_Configuration_: (optional) single map of options (`map[string]any`), provided as the single configuration argument in the rule's arguments array.

- `skip-convention-name-check`: (bool) If `true`, skip checks for package name conventions (underscores, MixedCaps, etc.). Default: `false`.
  This option is mutually exclusive with `convention-name-check-regex`; setting both results in a configuration error.
- `convention-name-check-regex`: (string) Custom regex pattern to validate package names. If set, package names must match this pattern.
  The value must be a non-empty string, and this option is mutually exclusive with `skip-convention-name-check`;
  setting both results in a configuration error.
- `skip-top-level-check`: (bool) If `true`, skip checks for top-level package names (e.g., `pkg`). Default: `false`.
- `skip-default-bad-name-check`: (bool) If `true`, skip checks for default bad package names (e.g., `common`, `utils`). Default: `false`.
- `check-extra-bad-name`: (bool) If `true`, enable checks for extra bad package names (e.g., `helpers`, `models`, `shared`, `utilities`). Default: `false`.
- `user-defined-bad-names`: ([]string) List of user-defined bad package names to check for.
- `skip-collision-with-common-std`: (bool) If `true`, skip checks for collisions with the most common Go standard library packages. Default: `false`.
- `check-collision-with-all-std`: (bool) If `true`, enable checks for collisions with all packages from Go standard library. Default: `false`.
  This option is mutually exclusive with `skip-collision-with-common-std`; setting both results in a configuration error.

Configuration examples:

Default settings (check for name conventions, top level packages, common bad names, and collision with common Go standard library packages):

```toml
[rule.package-naming]
```

Custom naming convention with regex:

```toml
[rule.package-naming]
arguments = [{ convention-name-check-regex = "^[a-z]+$" }]
```

Skip convention checks, but check for bad names:

```toml
[rule.package-naming]
arguments = [{ skip-convention-name-check = true }]
```

Enable collision checks with the most common standard library packages:

```toml
[rule.package-naming]
arguments = [{ skip-collision-with-common-std = false }]
```

Strict mode with user-defined bad names:

```toml
[rule.package-naming]
arguments = [{ user-defined-bad-names = ["foo", "bar"] }]
```

Enable collision checks with all standard library [packages](https://pkg.go.dev/std) excluding `internal` and `vendor`:

```toml
[rule.package-naming]
arguments = [{ check-collision-with-all-std = true }]
```

## package-directory-mismatch

_Go version_: 1.0.

_Revive version_: [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0).

_Description_: It is considered a good practice to name a package after the directory containing it.
This rule warns when the package name declared in the file does not match the name of the directory containing the file.

The following cases are excluded from this check:

- Package `main` (executable packages)
- Files in `testdata` directories (at any level) - by default
- Files directly in `internal` directories (but files in subdirectories of `internal` are checked)

For test files (files with `_test` suffix), package name additionally checked if it matches directory name with  `_test` suffix appended.

The rule normalizes both directory and package names before comparison by removing hyphens (`-`),
underscores (`_`), and dots (`.`). This allows package `foo_barbuz` be equal with directory `foo-bar.buz`.

For files in version directories (`v1`, `v2`, etc.), package name is checked if it matches either the version directory or its parent directory.

_Configuration_: Named arguments for directory exclusions.

Configuration examples:

Default behavior excludes paths containing `testdata`

```toml
[rule.package-directory-mismatch]
```

Ignore specific directories with `ignore-directories`

```toml
[rule.package-directory-mismatch]
arguments = [{ ignore-directories = ["testcases", "testinfo"] }]
```

Include all directories (`testdata` also)

```toml
[rule.package-directory-mismatch]
arguments = [{ ignore-directories = [] }]
```

## range-val-address

_Go version_: < 1.22.

_Revive version_: [v1.0.2](https://github.com/revive-lint/revive/releases/tag/v1.0.2).

**_Typed_**

_Description_: Range variables in a loop are reused at each iteration.
This rule warns when assigning the address of the variable, passing the address to append() or using it in a map.

_Configuration_: N/A

_Note_: This rule is irrelevant for Go 1.22+.

## range-val-in-closure

_Go version_: < 1.22.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Range variables in a loop are reused at each iteration; therefore a goroutine created in a loop will point to the range variable
with from the upper scope. This way, the goroutine could use the variable with an undesired value.
This rule warns when a range value (or index) is used inside a closure.

_Configuration_: N/A

_Note_: This rule is irrelevant for Go 1.22+.

## range

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: This rule suggests a shorter way of writing ranges that do not use the second value.

### Examples (range)

Before (violation):

```go
for i, _ := range items {
	process(i)
}
```

After (fixed):

```go
for i := range items {
	process(i)
}
```

_Configuration_: N/A

## receiver-naming

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: By convention, receiver names in a method should reflect their identity.
For example, if the receiver is of type `Parts`, `p` is an adequate name for it.
Contrary to other languages, it is not idiomatic to name receivers as `this` or `self`.
All methods of a type should also use the same receiver name.

### Examples (receiver-naming)

Before (violation):

```go
func (p *Parts) Add(part Part) {
	// ...
}

func (parts *Parts) Clear() {
	// ...
}
```

After (fixed):

```go
func (p *Parts) Add(part Part) {
	// ...
}

func (p *Parts) Clear() {
	// ...
}
```

_Configuration_: (optional) list of key-value-pair-map (`[]map[string]any`).

- `max-length`: (int) max length of receiver name

Configuration example:

```toml
[rule.receiver-naming]
arguments = [{ max-length = 2 }]
```

## redefines-builtin-id

_Go version_: 1.0; behavior changes in 1.21.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Constant names like `false`, `true`, `nil`, function names like `append`, `make`,
and basic type names like `bool`, and `byte` are not reserved words of the language; therefore the can be redefined.
Even if possible, redefining these built in names can lead to bugs very difficult to detect.

_Configuration_: N/A

## redundant-build-tag

_Go version_: >= 1.17; behavior changes in 1.21.

_Revive version_: [v1.6.0](https://github.com/revive-lint/revive/releases/tag/v1.6.0).

_Description_: This rule warns about redundant [build tag comments](https://pkg.go.dev/cmd/go@go1.17.0#hdr-Build_constraints).
It detects unnecessary `// +build` comments when `//go:build` is present.
`gofmt` in Go 1.17+ automatically adds the `//go:build` constraint, making the `// +build` comment unnecessary.
Also, the rule spots redundant build tags `//go:build go1.X` when the package's Go language version is greater than or equal to `go1.X`.

### Examples (redundant-build-tag)

Redundant `// +build` comment:

Before (violation):

```go
//go:build go1.17
// +build go1.17

package example
```

After (fixed):

```go
//go:build go1.17

package example
```

Redundant build tag when the module Go version is Go 1.21 or later:

Before (violation):

```go
//go:build go1.20

package example
```

After (fixed):

```go
package example
```

_Configuration_: N/A

_Note_: This rule is irrelevant for Go 1.16-.

## redundant-import-alias

_Go version_: 1.0.

_Revive version_: [v1.3.3](https://github.com/revive-lint/revive/releases/tag/v1.3.3).

_Description_: This rule warns on redundant import aliases. This happens when the alias used on the import statement matches the imported package name.

_Configuration_: N/A

## redundant-test-main-exit

_Go version_: >= 1.15.

_Revive version_: [v1.6.1](https://github.com/revive-lint/revive/releases/tag/v1.6.1).

_Description_: This rule warns about redundant `os.Exit` and `syscall.Exit` calls in the `TestMain` function,
as the Go test runner automatically handles program termination starting from Go 1.15.
Only calls that exit with the result of `m.Run` are reported (e.g. `os.Exit(m.Run())` or `code := m.Run(); os.Exit(code)`),
since exiting with another status is the only way to fail the test binary when setup or teardown fails.

_Configuration_: N/A

_Note_: This rule is irrelevant for Go 1.14-.

## string-format

_Go version_: 1.0.

_Revive version_: [v1.0.7](https://github.com/revive-lint/revive/releases/tag/v1.0.7).

_Description_: This rule allows you to configure a list of regular expressions that string literals in certain function calls are checked against.
This is geared towards user facing applications where string literals are often used for messages that will be presented to users,
so it may be desirable to enforce consistent formatting.

_Configuration_: Each argument is a slice containing 2-3 strings: a scope, a regex, and an optional error message.

1. The first string defines a **scope**. This controls which string literals the regex will apply to, and is defined as a function argument.
It must contain at least a function name (`core.WriteError`).
Scopes may optionally contain a number specifying which argument in the function to check (`core.WriteError[1]`),
as well as a struct field (`core.WriteError[1].Message`, only works for top level fields).
Function arguments are counted starting at 0, so `[0]` would refer to the first argument, `[1]` would refer to the second, etc.
If no argument number is provided, the first argument will be used (same as `[0]`).
You can use multiple scopes to one regex. Split them by `,` (`core.WriteError,fmt.Errorf`).

2. The second string is a **regular expression** (beginning and ending with a `/` character), which will be used to check the string literals in the scope.
The default semantics is "_strings matching the regular expression are OK_".
If you need to inverse the semantics you can add a `!` just before the first `/`. Examples:

    - with `"/^[A-Z].*$/"` the rule will **accept** strings starting with capital letters
    - with `"!/^[A-Z].*$/"` the rule will a **fail** on strings starting with capital letters

3. The third string (optional) is a **message** containing the purpose for the regex, which will be used in lint errors.

Configuration example:

```toml
[rule.string-format]
arguments = [
  [
    "core.WriteError[1].Message",
    "/^([^A-Z]|$)/",
    "must not start with a capital letter",
  ],
  [
    "fmt.Errorf[0]",
    "/(^|[^\\.!?])$/",
    "must not end in punctuation",
  ],
  [
    "panic",
    "/^[^\\n]*$/",
    "must not contain line breaks",
  ],
  [
    "fmt.Errorf[0],core.WriteError[1].Message",
    "!/^.*%w.*$/",
    "must not contain '%w'",
  ],
]
```

## string-of-int

_Go version_: 1.0.

_Revive version_: [v1.0.2](https://github.com/revive-lint/revive/releases/tag/v1.0.2).

**_Typed_**

_Description_: Explicit type conversion `string(i)` where `i` has an integer type other than `rune` might behave not as expected by the developer
(e.g. `string(42)` is not `"42"`). This rule spot that kind of suspicious conversions.

_Configuration_: N/A

## struct-tag

_Go version_: 1.0; behavior changes in 1.24.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: The rule spots errors in struct tags.
This is useful because struct tags are not checked at compile time.

The list of [supported tags](https://go.dev/wiki/Well-known-struct-tags):

| Tag           | Documentation                                                            |
| ------------- | ------------------------------------------------------------------------ |
| `asn1`         | <https://pkg.go.dev/encoding/asn1>                                      |
| `bson`         | <https://pkg.go.dev/go.mongodb.org/mongo-driver/bson>                   |
| `cbor`         | <https://pkg.go.dev/github.com/fxamacker/cbor/v2>                   |
| `datastore`    | <https://pkg.go.dev/cloud.google.com/go/datastore>                      |
| `default`      | The type of "default" must match the type of the field.                 |
| `json`         | <https://pkg.go.dev/encoding/json>                                      |
| `mapstructure` | <https://pkg.go.dev/github.com/mitchellh/mapstructure>                  |
| `properties`   | <https://pkg.go.dev/github.com/magiconair/properties#Properties.Decode> |
| `protobuf`     | <https://github.com/golang/protobuf>                                    |
| `required`     | Should be only "true" or "false".                                       |
| `spanner`      | <https://pkg.go.dev/cloud.google.com/go/spanner>                        |
| `toml`         | <https://pkg.go.dev/github.com/pelletier/go-toml/v2>                    |
| `url`          | <https://github.com/google/go-querystring>                              |
| `validate`     | <https://github.com/go-playground/validator>                            |
| `xml`          | <https://pkg.go.dev/encoding/xml>                                       |
| `yaml`         | <https://pkg.go.dev/gopkg.in/yaml.v2>                                   |

_Configuration_: (optional) The list of struct tags that can be accepted by the rule additionally to the supported tags.

Configuration example:

To accept the `inline` option in JSON tags (and `outline` and `gnu` in BSON tags) you must provide the following configuration

```toml
[rule.struct-tag]
arguments = ["json,inline", "bson,outline,gnu"]
```

To prevent a tag from being checked, simply add a `!` before its name.
For example, to instruct the rule not to check `validate` tags
(and accept `outline` and `gnu` in BSON tags) you can provide the following configuration

```toml
[rule.struct-tag]
arguments = ["!validate", "bson,outline,gnu"]
```

## superfluous-else

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: To improve the readability of code, it is recommended to reduce the indentation as much as possible.
This rule highlights redundant _else-blocks_ that can be eliminated when the preceding `if`-block deviates control flow,
for example ending with a `break`, `continue`, `goto`, `panic` or `os.Exit` call (the `return` case is handled by [indent-error-flow](#indent-error-flow)).

### Examples (superfluous-else)

Before (violation):

```go
for _, v := range values {
	if v < 0 {
		continue
	} else {
		sum += v
	}
}
```

After (fixed):

```go
for _, v := range values {
	if v < 0 {
		continue
	}

	sum += v
}
```

_Configuration_: ([]string) rule flags. Available flags are:

- `preserve-scope`: (string) do not suggest refactorings that would increase variable scope

Configuration example:

```toml
[rule.superfluous-else]
arguments = ["preserve-scope"]
```

## time-date

_Go version_: 1.0.

_Revive version_: [v1.10.0](https://github.com/revive-lint/revive/releases/tag/v1.10.0).

_Description_: Reports bad usage of `time.Date`.

_Configuration_: N/A

_Examples_:

- Invalid dates reporting:

  - 0 for the month or day argument
  - out of bounds argument for the month (12), day (31), hour (23), minute (59), or seconds (59)
  - an invalid date: 31st of June, 29th of February in 2023, ...

- Non-decimal integers are used as arguments

  This includes:

  - leading zero notation like using 00 for hours, minutes, and seconds.
  - octal notation 0o1, 0o0 that are often caused by using gofumpt on leading zero notation.
  - padding zeros such as 00123456 that are source of bugs.
  - ... and some other use cases.

<details>
<summary>More information about what is detected and reported</summary>

```go
import "time"

var (
	// here we can imagine zeroes were used for padding purpose
	a = time.Date(2023, 1, 2, 3, 4, 0, 00000000, time.UTC) // 00000000 is octal and equals 0 in decimal
	b = time.Date(2023, 1, 2, 3, 4, 0, 00000006, time.UTC) // 00000006 is octal and equals 6 in decimal
	c = time.Date(2023, 1, 2, 3, 4, 0, 00123456, time.UTC) // 00123456 is octal and equals 42798 in decimal
)
```

But here, what was expected 123456 or 42798 ? It's a source of bugs.

So the rule will report this

> octal notation with padding zeroes found: choose between 123456 and 42798 (decimal value of 123456 octal value)

This rule also reports strange notations used with time.Date.

Example:

```go
import "time"

var _ = time.Date(
	0x7e7,    // hexadecimal notation: use 2023 instead of 0x7e7/
	0b1,      // binary notation: use 1 instead of 0b1/
	0x_2,     // hexadecimal notation: use 2 instead of 0x_2/
	1_3,      // alternative notation: use 13 instead of 1_3/
	1e1,      // exponential notation: use 10 instead of 1e1/
	0.,       // float literal: use 0 instead of 0./
	0x1.Fp+6, // float literal: use 124 instead of 0x1.Fp+6/
	time.UTC)
```

All these are considered to be an uncommon usage of time.Date, are reported with a 0.8 confidence.

Note: even if 00, 01, 02, 03, 04, 05, 06, 07 are octal numbers, they can be considered as valid, and reported with 0.5 confidence.

```go
import "time"

var _ = time.Date(2023, 01, 02, 03, 04, 00, 0, time.UTC)
```

</details>

## time-equal

_Go version_: 1.0.

_Revive version_: [v1.1.3](https://github.com/revive-lint/revive/releases/tag/v1.1.3).

**_Typed_**

_Description_: This rule warns when using `==` and `!=` for equality checks on `time.Time` and suggests using the `time.Time.Equal` method.
For more information, see [this link](https://pkg.go.dev/time#Time).

_Configuration_: N/A

## time-naming

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

**_Typed_**

_Description_: Using unit-specific suffix like "Secs", "Mins", ... when naming variables of type `time.Duration` can be misleading,
this rule highlights those cases.

_Configuration_: N/A

## unchecked-type-assertion

_Go version_: 1.0.

_Revive version_: [v1.3.4](https://github.com/revive-lint/revive/releases/tag/v1.3.4).

_Description_: This rule checks whether a type assertion result is checked (the `ok` value), preventing unexpected `panic`s.

_Configuration_: list of key-value-pair-map (`[]map[string]any`).

- `accept-ignored-assertion-result`: (bool) default `false`,
set it to `true` to accept ignored type assertion results like this:

```golang
foo, _ := bar(.*Baz).
//   ^
```

Configuration example:

```toml
[rule.unchecked-type-assertion]
arguments = [{ accept-ignored-assertion-result = true }]
```

## unconditional-recursion

_Go version_: 1.0.

_Revive version_: [v1.0.3](https://github.com/revive-lint/revive/releases/tag/v1.0.3).

_Description_: Unconditional recursive calls will produce infinite recursion, thus program stack overflow.
This rule detects and warns about unconditional (direct) recursive calls.

_Configuration_: N/A

## unexported-naming

_Go version_: 1.0.

_Revive version_: [v1.0.3](https://github.com/revive-lint/revive/releases/tag/v1.0.3).

_Description_: this rule warns on wrongly named un-exported symbols, i.e. un-exported symbols whose name start with a capital letter.

_Configuration_: N/A

## unexported-return

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

**_Typed_**

_Description_: This rule warns when an exported function or method returns a value of an un-exported type.
Unexported interface types are reported too.
To expose an unexported type, declare an exported alias for it (e.g. `type Option = option`) and return the alias.

_Configuration_: N/A

## unhandled-error

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Typed_**

_Description_: This rule warns when errors returned by a function are not explicitly handled on the caller side.

_Configuration_: function names regexp patterns to ignore

Configuration example:

```toml
[rule.unhandled-error]
arguments = [
  '^os\.(CreateTemp|WriteFile|Chmod)$',
  '^fmt\.Print',
  'myFunction',
  '^net\.',
  '^(bytes\.Buffer|string\.Writer)\.Write(Byte|Rune|String)?$',
]
```

## unnecessary-if

_Go version_: 1.0.

_Revive version_: [v1.13.0](https://github.com/revive-lint/revive/releases/tag/v1.13.0).

_Description_: Detects unnecessary `if-else` statements that return or assign a boolean value
based on a condition and suggests a simplified, direct return or assignment.
The `if-else` block is redundant because the condition itself is already a boolean expression.
The simplified version is immediately clearer, more idiomatic, and reduces cognitive load for the reader.

### Examples (unnecessary-if)

```go
if y <= 0 {
  z = true
} else {
  z = false
}

if x > 10 {
  return false
} else {
  return true
}
```

Fixed code:

```go
z = y <= 0

return x <= 10
```

_Configuration_: N/A

## unnecessary-format

_Go version_: 1.0.

_Revive version_: [v1.10.0](https://github.com/revive-lint/revive/releases/tag/v1.10.0).

_Description_: This rule identifies calls to formatting functions where the format string does not contain any formatting verbs
and recommends switching to the non-formatting, more efficient alternative.

_Configuration_: N/A

## unnecessary-stmt

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: This rule suggests to remove redundant statements like a `break` at the end of a case block, for improving the code's readability.

### Examples (unnecessary-stmt)

Before (violation):

```go
switch status {
case "active":
	handle()
	break
}
```

After (fixed):

```go
switch status {
case "active":
	handle()
}
```

_Configuration_: N/A

## unreachable-code

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: This rule spots and proposes to remove [unreachable code](https://en.wikipedia.org/wiki/Unreachable_code).

### Examples (unreachable-code)

Before (violation):

```go
import "log"

func compute() int {
	return doWork()
	log.Println("done")
}
```

After (fixed):

```go
import "log"

func compute() int {
	log.Println("starting")
	return doWork()
}
```

_Configuration_: N/A

## unsecure-url-scheme

_Go version_: 1.0.

_Revive version_: [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0).

_Description_: Checks for usage of potentially unsecure URL schemes (`http`, `ws`) in string literals.
Using unencrypted URL schemes can expose sensitive data during transmission and
make applications vulnerable to man-in-the-middle attacks.
Secure alternatives like `https` should be preferred when possible.

_Configuration_: N/A

The rule will not warn on local URLs (`localhost`, `127.0.0.1`).

## unused-parameter

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: This rule warns on unused parameters. Functions or methods with unused parameters can be a symptom of an unfinished refactoring or a bug.

_Configuration_: Supports a single `map[string]any` argument with an `allow-regex` option to
specify additional allowed patterns for unused parameter names beyond the default `_`.

Configuration example:

This allows any names starting with `_`, not just `_` itself:

```go
func SomeFunc(_someObj *MyStruct) {} // matches rule
```

```toml
[rule.unused-parameter]
arguments = [{ allow-regex = "^_" }]
```

## unused-receiver

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: This rule warns on unused method receivers. Methods with unused receivers can be a symptom of an unfinished refactoring or a bug.

_Configuration_:
Supports a single `map[string]any` argument with an `allow-regex` option to specify additional allowed unused receiver name patterns beyond `_`.

Configuration example:

This allows any names starting with `_`, not just `_` itself:

```go
func (_my *MyStruct) SomeMethod() {} // matches rule
```

```toml
[rule.unused-receiver]
arguments = [{ allow-regex = "^_" }]
```

## use-any

_Go version_: >= 1.18.

_Revive version_: [v1.2.0](https://github.com/revive-lint/revive/releases/tag/v1.2.0).

_Description_: This rule proposes to replace instances of `interface{}` with its alias [`any`](https://pkg.go.dev/builtin@go1.18.0#any).

### Examples (use-any)

Before (violation):

```go
import "fmt"

func PrintValue(v interface{}) {
	fmt.Println(v)
}
```

After (fixed):

```go
import "fmt"

func PrintValue(v any) {
	fmt.Println(v)
}
```

_Configuration_: N/A

_Note_: This rule is irrelevant for Go 1.17-.

## use-errors-new

_Go version_: < 1.26.

_Revive version_: [v1.6.0](https://github.com/revive-lint/revive/releases/tag/v1.6.0).

_Description_: This rule identifies calls to `fmt.Errorf` that can be safely replaced by, the more efficient, `errors.New`.
This applies when the format string has no formatting verbs (no additional arguments are passed).

### Examples (use-errors-new)

Before (violation):

```go
import "fmt"

return fmt.Errorf("connection refused")
```

After (fixed):

```go
import "errors"

return errors.New("connection refused")
```

_Configuration_: N/A

_Note_: This rule is irrelevant for Go 1.26+.
For unformatted strings, `fmt.Errorf("x")` generally [matches](https://go.dev/doc/go1.26#fmtpkgfmt) the allocations for `errors.New("x")`.

## use-fmt-print

_Go version_: 1.0.

_Revive version_: [v1.10.0](https://github.com/revive-lint/revive/releases/tag/v1.10.0).

_Description_: This rule proposes to replace calls to built-in `print` and `println` with their equivalents from `fmt` standard package.

`print` and `println` built-in functions are not recommended for use-cases other than
[language bootstrapping and are not guaranteed to stay in the language](https://go.dev/ref/spec#Bootstrapping).

_Configuration_: N/A

## use-slices-concat

_Go version_: >= 1.22.

_Revive version_: [v1.17.0](https://github.com/revive-lint/revive/releases/tag/v1.17.0).

_Description_: Since Go 1.22 the `slices` package provides the `slices.Concat` function that concatenates slices into a new one.
The rule proposes to replace appends to an empty slice, which do the same in a more convoluted way, with a call to `slices.Concat`.

### Examples (use-slices-concat)

Before (violation):

```go
content := append(append([]byte{}, magic...), payload...)

allRanges := append([]ignoredRange{}, inlineRanges...)
allRanges = append(allRanges, expandedRanges...)
```

After (fixed):

```go
content := slices.Concat(magic, payload)

allRanges := slices.Concat(inlineRanges, expandedRanges)
```

_Configuration_: N/A

_Note_: The rule does not use type information, thus it can report appends that `slices.Concat` can not replace,
for example the append of a `string` to a `[]byte`.
Mind also that `slices.Concat` returns `nil` where the appends return an empty slice when all the concatenated slices are empty.

_Note_: This rule is irrelevant for Go 1.21-.

## use-slices-sort

_Go version_: >= 1.21.

_Revive version_: [v1.14.0](https://github.com/revive-lint/revive/releases/tag/v1.14.0).

_Description_: Since Go 1.21 the `slices` package proposes methods that are faster and easier to use
than their equivalents in `sort` package.
The rule proposes to replace these legacy idioms with calls to the new methods.

### Examples (use-slices-sort)

```go
sort.Float64s(temperatures)
sort.Ints(years)
sort.Strings(names)
sort.Slice(cfg.Dependencies, func(i, j int) bool {
	return cfg.Dependencies[i].Name < cfg.Dependencies[j].Name
})
```

Fixed code:

```go
slices.Sort(temperatures)
slices.Sort(years)
slices.Sort(names)
slices.SortFunc(cfg.Dependencies, func(a, b config.Dependency) int {
	return cmp.Compare(a.Name, b.Name)
})
```

_Configuration_: N/A

_Note_: This rule is irrelevant for Go 1.20-.

## use-waitgroup-go

_Go version_: >= 1.25.

_Revive version_: [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0).

_Description_: Since Go 1.25 the `sync` package proposes the [`WaitGroup.Go`](https://pkg.go.dev/sync#WaitGroup.Go) method.
This method is a shorter and safer replacement for the idiom `wg.Add ... go { ... wg.Done ... }`.
The rule proposes to replace these legacy idioms with calls to the new method.

_Limitations_: The rule doesn't rely on type information but on variable names to identify waitgroups.
This means the rule search for `wg` (the defacto standard name for wait groups);
if the waitgroup variable is named differently than `wg` the rule will skip it.

### Examples (use-waitgroup-go)

Before (violation):

```go
import "sync"

func concurrentProcess(jobs int) {
	wg := sync.WaitGroup{}

	wg.Add(jobs)

	for range jobs {
		go func() {
			defer wg.Done()
			// do something
		}()
	}

	wg.Wait()
}
```

After (fixed):

```go
import "sync"

func concurrentProcess(jobs int) {
	wg := sync.WaitGroup{}

	for range jobs {
		wg.Go(func() {
			// do something
		})
	}

	wg.Wait()
}
```

_Configuration_: N/A

_Note_: This rule is irrelevant for Go 1.24-.

## useless-break

_Go version_: 1.0.

_Revive version_: [v1.1.0](https://github.com/revive-lint/revive/releases/tag/v1.1.0).

_Description_: This rule warns on useless `break` statements in case clauses of switch and select statements. Go,
unlike other programming languages like C, only executes statements of the selected case while ignoring the subsequent case clauses.
Therefore, inserting a `break` at the end of a case clause has no effect.

Because `break` statements are rarely used in case clauses, when switch or select statements are inside a for-loop,
the programmer might wrongly assume that a `break` in a case clause will take the control out of the loop.
The rule emits a specific warning for such cases.

### Examples (useless-break)

Before (violation):

```go
for {
	switch state {
	case done:
		cleanup()
		break
	}
}
```

After (fixed):

```go
for {
	switch state {
	case done:
		cleanup()
	}
}
```

_Configuration_: N/A

## useless-fallthrough

_Go version_: 1.0.

_Revive version_: [v1.12.0](https://github.com/revive-lint/revive/releases/tag/v1.12.0).

_Description_: This rule warns on useless `fallthrough` statements in case clauses of switch statements.
A `fallthrough` is considered _useless_ if it's the single statement of a case clause block.

Go allows `switch` statements with clauses that group multiple cases.
Thus, for example:

```go
switch category {
  case "Lu":
    fallthrough
  case "Ll":
    fallthrough    
  case "Lt":
    fallthrough
  case "Lm": 
    fallthrough
  case "Lo":
    return true
  default:
    return false
}
```

can be written as

```go
switch category {
  case "Lu", "Ll", "Lt", "Lm", "Lo":
      return true
  default:
    return false
}
```

_Configuration_: N/A

## var-declaration

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

**_Typed_**

_Description_: This rule proposes simplifications of variable declarations.

_Configuration_: N/A

## var-naming

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

**_Ported from golint_**

_Description_: This rule warns when [initialism](https://go.dev/wiki/CodeReviewComments#initialisms), [variable](https://go.dev/wiki/CodeReviewComments#variable-names)
naming conventions are not followed.
It ignores functions starting with `Example`, `Test`, `Benchmark`, and `Fuzz` in test files, preserving `golint` original behavior.

_Configuration_: This rule accepts two slices of strings and one optional slice containing a single map with named parameters.
(This is because TOML does not support "slice of any," and we maintain backward compatibility with the previous configuration version).
The first slice is an allowlist, and the second one is a blocklist of initialisms.
You can add a boolean parameter `skip-initialism-name-checks` to control how names
of functions, variables, consts, and structs handle known initialisms (e.g., JSON, HTTP, etc.) when written in `camelCase`.
When `skip-initialism-name-checks` is set to true, the rule allows names like `readJson`, `HttpMethod` etc.
You can add a boolean parameter `initialisms-as-words` to treat initialisms as normal words
for uppercase and lowercase rules (forbidding names like `readJSON`, `HTTPMethod`, `ID` and requiring `readJson`, `HttpMethod`, `Id`, etc.).
When `initialisms-as-words` is set to true, the built-in list of common initialisms is ignored,
while allowlist and blocklist configurations remain active.
In the map, you can add a boolean `upper-case-const` parameter to allow `UPPER_CASE` for `const`.

By default, the rule behaves exactly as the alternative in `golint` for non-package identifiers;
`golint`-equivalent package-name warnings now require enabling the [`package-naming`](#package-naming) rule.
The legacy package-related options `skip-package-name-checks`, `extra-bad-package-names`, and `skip-package-name-collision-with-go-std` are deprecated
and are now treated as no-ops by `var-naming` (they are ignored, apart from an optional warning when logging is enabled).
Package-name checks should be configured via the [`package-naming`](#package-naming) rule instead,
and these options should be removed from `var-naming` configurations to avoid confusion.

Configuration examples:

```toml
[rule.var-naming]
arguments = [[], [], [{ skip-initialism-name-checks = true }]]
```

```toml
[rule.var-naming]
arguments = [[], [], [{ initialisms-as-words = true }]]
```

```toml
[rule.var-naming]
arguments = [["ID"], ["VM"], [{ upper-case-const = true }]]
```

## waitgroup-by-value

_Go version_: 1.0.

_Revive version_: [v1.0.0](https://github.com/revive-lint/revive/releases/tag/v1.0.0).

_Description_: Function parameters that are passed by value, are in fact a copy of the original argument.
Passing a copy of a `sync.WaitGroup` is usually not what the developer wants to do.
This rule warns when a `sync.WaitGroup` expected as a by-value parameter in a function or method.

_Configuration_: N/A
