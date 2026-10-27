package fixtures // MATCH /the file has a comment density of 43% (6 comment lines for 8 code lines) but expected a minimum of 60%/

/*
Block with empty line

Line 4
*/
func sample() {
	fakeCall(
		"// not comment",
		`/* also
not comment */`,
	)
}
