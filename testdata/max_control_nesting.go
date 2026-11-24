package fixtures

func mcn() {
	if true {
		if true {
			if true { // MATCH /control flow nesting exceeds 2/

			}
		}
	} else {
		if true {
			if true { // MATCH /control flow nesting exceeds 2/
				if true {

				}
			}
		}
	}

	for {
		if true {
			for { // MATCH /control flow nesting exceeds 2/
			}
		}
	}

	for range xs {
		if true {
			for range xs { // MATCH /control flow nesting exceeds 2/
			}
		}
	}

	switch {
	case false:
		if true {

		}
	case true:
		if true {
			for { // MATCH /control flow nesting exceeds 2/
			}
		}
	default:
	}

	select {
	case msg1 := <-c1:
		println("received", msg1)
	case msg2 := <-c2:
		if true {
			for { // MATCH /control flow nesting exceeds 2/
			}
		}
	}

	if true {
		f1 := func() {
			if true {
				for {
				}
			}
		}
	}

	f1 := func() {
		for {
			if true {
				for { // MATCH /control flow nesting exceeds 2/
				}
			}
		}
	}
}
