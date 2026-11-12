// Package goreleaser is a placeholder for tooling imports, and should not be imported in production code.
// Despite `go tool` this package is needed to workaround the fact that renovate can't update indirect dependencies.
package goreleaser

import (
	_ "github.com/goreleaser/goreleaser/v2/pkg/config" //revive:disable-line:blank-imports
)
