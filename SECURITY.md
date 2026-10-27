# Security Policy

## Supported versions

Only the latest release is supported. Fixes are cut from `master` as a new release; there are no release branches and no backports.
See [Releasing](DEVELOPING.md#releasing) in the developer's guide.

## Reporting a vulnerability

Report vulnerabilities privately through GitHub Security Advisories:

**<https://github.com/revive-lint/revive/security/advisories/new>**

The report stays visible only to you and the maintainers until an advisory is published.
Please do not open a public issue or pull request for a vulnerability.

Include whatever you have: the affected version, what an attacker gains, and the smallest input or configuration that reproduces it.

If GitHub Security Advisories is not an option for you, contact a [maintainer](README.md#maintainers) privately instead.

## What to expect

- An initial response within 14 days, acknowledging the report and saying whether we consider it a vulnerability.
- Updates as we work on a fix, and credit in the advisory and the release notes unless you prefer otherwise.
- A fix released as a patch version once it is ready, with a published advisory.

## Scope

revive is a linter: it reads Go source files and writes findings to standard output or a file.
It opens no network connections, implements no cryptography, and handles no credentials.
Reports about the tool itself are in scope — for example arbitrary code execution or file access triggered by linting untrusted source,
by a configuration file, or by a custom rule loaded through the library API.

Findings that revive reports, or fails to report, in the code it lints are not vulnerabilities in revive.
A false positive or false negative in a rule, including a security-related one, is a [bug report](https://github.com/revive-lint/revive/issues/new/choose).
