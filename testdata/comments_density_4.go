package fixtures // MATCH /the file has a comment density of 50% (5 comment lines for 5 code lines) but expected a minimum of 60%/

// Package level declaration
type Options struct {
	/* Block comment
	with a second line
	and a third line */
	Timeout int
}

var url = "http://example.com"
