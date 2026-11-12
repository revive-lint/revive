// Package golangcilint is a placeholder for tooling imports, and should not be imported in production code.
// Despite `go tool` this package is needed to workaround the fact that renovate can't update indirect dependencies.
package golangcilint

import (
	_ "github.com/golangci/golangci-lint/v2/pkg/exitcodes" //revive:disable-line:blank-imports
)
