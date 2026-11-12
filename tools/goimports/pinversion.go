// Package goimports is a placeholder for tooling imports, and should not be imported in production code.
// Despite `go tool` this package is needed to workaround the fact that renovate can't update indirect dependencies.
package goimports

import (
	_ "golang.org/x/tools/imports" //revive:disable-line:blank-imports
)
