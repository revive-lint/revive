//go:build linux && go1.21

// MATCH:1 /The build tag "//go:build linux && go1.21" is implied by the file name "redundant_build_tag_version_linux.go" and the go directive in go.mod, so it can be removed/

package fixtures
