//go:build linux
// +build linux

// MATCH:2 /The build tag "// +build" is redundant since Go 1.17 and can be removed/

// The `// +build` tag is reported even though the `//go:build` one is implied by the file name:
// removing the `//go:build` line would only make gofmt add it back from the `// +build` one.

package pkg
