//go:build darwin && cgo

// MATCH:1 /The build tag "//go:build darwin && cgo" is partly implied by the file name "redundant_build_tag_partly_darwin.go" and can be reduced to "//go:build cgo"/

package pkg
