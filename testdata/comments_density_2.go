package fixtures // MATCH /the file has a comment density of 15% (5 comment lines for 29 code lines) but expected a minimum of 90%/

// datarace is function
func datarace() (r int, c char) {
	for _, p := range []int{1, 2} {
		go func() {
			print(r)
			print(p)
		}()
		for i, p1 := range []int{1, 2} {
			a := p1
			go func() {
				print(r)
				print(p)
				print(p1)
				print(a)
				print(i)
			}()
			print(i)
			print(p)
			go func() {
				_ = c
			}()
		}
		print(p1)
	}
	/* Goroutines
	are
	awesome */
	go func() {
		print(r)
	}()
	print(r)
}
