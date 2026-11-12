package pkg

// Reported: the variable exists only to have its address taken.

func assignToField(size *float32, values []float64, i int) {
	sizeVal := float64(*size)
	snapshot.Size = &sizeVal // MATCH /use new(float64(*size)) instead of the temporary variable sizeVal/

	nameVal := "a name"
	snapshot.Name = &nameVal // MATCH /use new("a name") instead of the temporary variable nameVal/

	timeout := config.Timeout
	snapshot.Size = &timeout // MATCH /use new(config.Timeout) instead of the temporary variable timeout/

	item := values[i]
	snapshot.Size = &item // MATCH /use new(values[i]) instead of the temporary variable item/
}

func definePointer() *int {
	count := 42
	countPtr := &count // MATCH /use new(42) instead of the temporary variable count/

	return countPtr
}

func returnAddress(b []byte) *string {
	s := string(b)

	return &s // MATCH /use new(string(b)) instead of the temporary variable s/
}

func passAddress(flag bool) {
	enabled := !flag
	consume(&enabled) // MATCH /use new(!flag) instead of the temporary variable enabled/
}

func fillCompositeLiteral(b []byte) *snapshotType {
	name := string(b)

	return &snapshotType{Name: &name} // MATCH /use new(string(b)) instead of the temporary variable name/
}

func inClauses(ch chan int, i int) {
	switch i {
	case 0:
		inCase := int64(i)
		int64s = append(int64s, &inCase) // MATCH /use new(int64(i)) instead of the temporary variable inCase/
	}

	select {
	case v := <-ch:
		inComm := int64(v)
		int64s = append(int64s, &inComm) // MATCH /use new(int64(v)) instead of the temporary variable inComm/
	}
}

func sameNameAsField(i int) {
	size := float64(i)
	snapshot.Size = &size // MATCH /use new(float64(i)) instead of the temporary variable size/
	other.size = 2
	consumeSnapshot(&snapshotType{size: 3})
}

func chainedPointers() **int {
	count := 42
	countPtr := &count // MATCH /use new(42) instead of the temporary variable count/

	return &countPtr
}

func sendAndDefer(ch chan *int64, i int) {
	sent := int64(i)
	ch <- &sent // MATCH /use new(int64(i)) instead of the temporary variable sent/

	deferred := int64(i)
	defer consume(&deferred) // MATCH /use new(int64(i)) instead of the temporary variable deferred/
}

func besideSideEffectFreeCalls(i int) *snapshotType {
	confidence := 0.9
	return &snapshotType{Confidence: &confidence, Size: new(0.5), Count: new(int64(i))} // MATCH /use new(0.9) instead of the temporary variable confidence/
}

func literalBesideCalls(i int) *snapshotType {
	name := "a name"
	consumeTwo(lookupSize("a name"), &snapshotType{Name: &name}) // MATCH /use new("a name") instead of the temporary variable name/

	negative := float64(-1)
	return newBuilder("a name").withSize(&negative) // MATCH /use new(float64(-1)) instead of the temporary variable negative/
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
