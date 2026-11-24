package fixtures

func unnecessaryIf() bool {
	var cond bool
	var id bool
	var other bool

	// test return replacements
	if cond { // MATCH /replace this conditional by: return cond/
		return true
	} else {
		return false
	}

	if cond { // MATCH /replace this conditional by: return !(cond)/
		return false
	} else {
		return true
	}

	// test assignment replacements
	if cond { // MATCH /replace this conditional by: id = cond/
		id = true
	} else {
		id = false
	}

	if cond { // MATCH /replace this conditional by: id = !(cond)/
		id = false
	} else {
		id = true
	}

	// test suggestions for (in)equalities
	//// assignments
	if cond == id { // MATCH /replace this conditional by: id = cond == id/
		id = true
	} else {
		id = false
	}

	if cond == id { // MATCH /replace this conditional by: id = cond != id/
		id = false
	} else {
		id = true
	}

	if cond != id { // MATCH /replace this conditional by: id = cond != id/
		id = true
	} else {
		id = false
	}

	if cond != id { // MATCH /replace this conditional by: id = cond == id/
		id = false
	} else {
		id = true
	}

	//// return
	if cond == id { // MATCH /replace this conditional by: return cond == id/
		return true
	} else {
		return false
	}

	if cond == id { // MATCH /replace this conditional by: return cond != id/
		return false
	} else {
		return true
	}

	if cond != id { // MATCH /replace this conditional by: return cond != id/
		return true
	} else {
		return false
	}

	if cond != id { // MATCH /replace this conditional by: return cond == id/
		return false
	} else {
		return true
	}

	//// assignments
	if cond <= id { // MATCH /replace this conditional by: id = cond <= id/
		id = true
	} else {
		id = false
	}

	if cond <= id { // MATCH /replace this conditional by: id = cond > id/
		id = false
	} else {
		id = true
	}

	if cond >= id { // MATCH /replace this conditional by: id = cond >= id/
		id = true
	} else {
		id = false
	}

	if cond >= id { // MATCH /replace this conditional by: id = cond < id/
		id = false
	} else {
		id = true
	}

	if cond > id { // MATCH /replace this conditional by: id = cond > id/
		id = true
	} else {
		id = false
	}

	if cond > id { // MATCH /replace this conditional by: id = cond <= id/
		id = false
	} else {
		id = true
	}

	if cond < id { // MATCH /replace this conditional by: id = cond < id/
		id = true
	} else {
		id = false
	}

	if cond < id { // MATCH /replace this conditional by: id = cond >= id/
		id = false
	} else {
		id = true
	}

	if (something > 0) && (!id) || (something+10 <= 0) { // MATCH /replace this conditional by: id = !((something > 0) && (!id) || (something+10 <= 0))/
		id = false
	} else {
		id = true
	}

	// conditionals with initialization
	if cond := false; cond {
		return true
	} else {
		return false
	}

	if cond := false; cond {
		id = true
	} else {
		id = false
	}

	// same value in both branches: no simplification possible
	if cond {
		return true
	} else {
		return true
	}

	if cond {
		return false
	} else {
		return false
	}

	if cond {
		id = true
	} else {
		id = true
	}

	if cond {
		id = false
	} else {
		id = false
	}

	// different assignment tokens: := and = must not be conflated
	if cond {
		other := true
	} else {
		other = false
	}

	return id == id
}
