// Package shfmt is a placeholder for tooling imports, and should not be imported in production code.
// Despite `go tool` this package is needed to workaround the fact that renovate can't update indirect dependencies.
package shfmt

import (
	_ "mvdan.cc/sh/v3/syntax" //revive:disable-line:blank-imports
)
