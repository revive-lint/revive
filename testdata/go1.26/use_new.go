package pkg

// Reported: the variable exists only to have its address taken.

func assignToField(size *float32, values []float64, i int) {
	sizeVal := float64(*size) // MATCH /replace the temporary variable sizeVal with new(float64(*size))/
	snapshot.Size = &sizeVal

	nameVal := "a name" // MATCH /replace the temporary variable nameVal with new("a name")/
	snapshot.Name = &nameVal

	timeout := config.Timeout // MATCH /replace the temporary variable timeout with new(config.Timeout)/
	snapshot.Size = &timeout

	item := values[i] // MATCH /replace the temporary variable item with new(values[i])/
	snapshot.Size = &item
}

func definePointer() *int {
	count := 42 // MATCH /replace the temporary variable count with new(42)/
	countPtr := &count

	return countPtr
}

func returnAddress(b []byte) *string {
	s := string(b) // MATCH /replace the temporary variable s with new(string(b))/

	return &s
}

func passAddress(flag bool) {
	enabled := !flag // MATCH /replace the temporary variable enabled with new(!flag)/
	consume(&enabled)
}

func fillCompositeLiteral(b []byte) *snapshotType {
	name := string(b) // MATCH /replace the temporary variable name with new(string(b))/

	return &snapshotType{Name: &name}
}

func inClauses(ch chan int, i int) {
	switch i {
	case 0:
		inCase := int64(i) // MATCH /replace the temporary variable inCase with new(int64(i))/
		int64s = append(int64s, &inCase)
	}

	select {
	case v := <-ch:
		inComm := int64(v) // MATCH /replace the temporary variable inComm with new(int64(v))/
		int64s = append(int64s, &inComm)
	}
}

func sameNameAsField(i int) {
	size := float64(i) // MATCH /replace the temporary variable size with new(float64(i))/
	snapshot.Size = &size
	other.size = 2
	consumeSnapshot(&snapshotType{size: 3})
}

func chainedPointers() **int {
	count := 42 // MATCH /replace the temporary variable count with new(42)/
	countPtr := &count

	return &countPtr
}

func sendAndDefer(ch chan *int64, i int) {
	sent := int64(i) // MATCH /replace the temporary variable sent with new(int64(i))/
	ch <- &sent

	deferred := int64(i) // MATCH /replace the temporary variable deferred with new(int64(i))/
	defer consume(&deferred)
}

// Not reported: the variable has another use, or is not declared right before its address is taken.

func otherUses(flag bool, i int, sizes [3]float64) {
	usedTwice := float64(i)
	snapshot.Size = &usedTwice
	_ = usedTwice

	var declared = float64(i)
	snapshot.Size = &declared

	multiValue, _ := lookup(i)
	snapshot.Size = &multiValue

	notAdjacent := float64(i)
	_ = flag
	snapshot.Size = &notAdjacent

	addressTakenTwice := float64(i)
	snapshot.Size = &addressTakenTwice
	snapshot.Size = &addressTakenTwice

	inFuncLit := !flag
	run(func() { consume(&inFuncLit) })

	mapKey := float64(i)
	snapshot.Size = &mapKey
	_ = map[float64]int{mapKey: 1}

	nestedMapKey := float64(i)
	snapshot.Size = &nestedMapKey
	_ = map[string]map[float64]int{"a": {nestedMapKey: 1}}

	arrayKey := sizes
	snapshot.Size = &arrayKey
	_ = [...]int{len(arrayKey): 1}
}

// Not reported: the address is taken in a body that can run many times, or not at all.

func inBody(values []float64, i int, n int, flag bool) {
	count := 0
	for range n {
		int64s = append(int64s, &count)
	}

	item := values[i]
	if flag {
		snapshot.Size = &item
	}

	inSwitch := float64(n)
	switch {
	case flag:
		snapshot.Size = &inSwitch
	}
}

// Not reported: evaluating the value where the address is taken would run it out of order.

func blockingValue(ch chan int64) {
	fromChannel := <-ch
	int64s = append(int64s, &fromChannel)

	alsoFromChannel := <-ch
	consume(&alsoFromChannel)
}

func besideAnotherOperation(ch chan int64, values []float64, i int) {
	beforeIndex := values[i]
	sizes[nextIndex()] = &beforeIndex

	beforeArgument := float64(i)
	consumeTwo(lookupSize("a name"), &beforeArgument)

	beforeReceive := float64(i)
	consumeTwo(<-ch, &beforeReceive)

	beforeCallee := float64(i)
	newBuilder("a name").consumeSize(&beforeCallee)
}

// Not reported: inlining the value would make the code harder to read.

func complexValues(name string) {
	compound := count(name) + 1
	snapshot.Size = &compound

	callResult := lookupSize(name)
	snapshot.Size = &callResult

	dereferenced := *newBuilder(name).Obj()
	consumeSnapshot(&dereferenced)

	literal := snapshotType{Name: snapshot.Name}
	consumeSnapshot(&literal)
}
