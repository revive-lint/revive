//go:build linux || windows

// MATCH:1 /The build tag "//go:build linux || windows" is implied by the file name "redundant_build_tag_or_linux.go" and can be removed/

package pkg
