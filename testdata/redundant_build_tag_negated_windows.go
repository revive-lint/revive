//go:build windows && !no_wsl

// MATCH:1 /The build tag "//go:build windows && !no_wsl" is partly implied by the file name "redundant_build_tag_negated_windows.go" and can be reduced to "//go:build !no_wsl"/

package pkg
