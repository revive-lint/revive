//go:build solaris

// Only GOOS=illumos satisfies the "solaris" tag, and `-tags=illumos` builds the file for any other GOOS,
// so the tag is not implied by the file name.

package pkg
