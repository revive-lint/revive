package fixtures

import "fmt"

type namedIntMap map[int]string

func inefficientMapLookup() {
	type aS struct {
		TagIDs map[int]string
	}
	a := aS{}
	someStaticValue := 1
	// use case from issue #1447
	for id := range a.TagIDs { // MATCH /inefficient lookup of map key/
		if id == someStaticValue {
			return
		}
	}

	var namedMap namedIntMap
	for key := range namedMap { // MATCH /inefficient lookup of map key/
		if key == someStaticValue {
			return
		}
	}

	for key, _ := range a.TagIDs { // MATCH /inefficient lookup of map key/
		if key == someStaticValue {
			return
		}
	}

	for key, _ := range a.TagIDs { // MATCH /inefficient lookup of map key/
		if key != someStaticValue {
			continue
		}
		fmt.Println(key)
	}

	// do not match if the loop body contains more than
	// just an if statement on the map key
	aMap := map[int]int{}
	for k := range aMap {
		fmt.Println(k)
		if k == 1 {
			return
		}
	}

	for k := range aMap {
		if k == 1 {
			return
		}
		fmt.Println(k)
	}

	// do not match on ranges over types other than maps
	slice := []int{}
	for i, _ := range slice {
		if i == 1 {
			fmt.Print(i)
		}
	}

	for key, _ := range a.TagIDs {
		if key != someStaticValue { // do not match if the loop body does more than just continuing
			fmt.Println(key)
			continue
		}
		fmt.Println(key)
	}

	// Test case for issue #1601
	var count int
	for range aMap {
		count++
	}
	_ = count
}
