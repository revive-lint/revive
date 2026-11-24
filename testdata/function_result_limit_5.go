package fixtures

func sixResults() (int, int, int, int, int, int) { // MATCH /maximum number of return results per function exceeded; max 5 but got 6/
	return 0, 0, 0, 0, 0, 0
}

func fiveResults() (int, int, int, int, int) {
	return 0, 0, 0, 0, 0
}
