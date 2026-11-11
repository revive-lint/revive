package fixtures

func Range(xs [][][][][]int) int {
	n := 0
	for _, a := range xs {
		for _, b := range a {
			for _, c := range b {
				for _, d := range c { // MATCH /control flow nesting exceeds 3/
					for _, e := range d {
						if e > 0 {
							n++
						}
					}
				}
			}
		}
	}
	return n
}
