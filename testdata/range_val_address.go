package fixtures

func rangeValAddress() {
	m := map[string]*string{}

	mySlice := []string{"A", "B", "C"}
	for _, value := range mySlice {
		m["address"] = &value // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}

func rangeValAddress2() {
	m := map[string]*string{}

	mySlice := []string{"A", "B", "C"}
	for i := range mySlice {
		m["address"] = &mySlice[i]
	}
}

func rangeValAddress3() {
	m := map[string]*string{}

	mySlice := []string{"A", "B", "C"}
	for _, value := range mySlice {
		a := &value // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
		m["address"] = a
	}
}

func rangeValAddress4() {
	m := []*string{}

	mySlice := []string{"A", "B", "C"}
	for _, value := range mySlice {
		m = append(m, &value) // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}

func rangeValAddress5() {
	m := map[*string]string{}

	mySlice := []string{"A", "B", "C"}
	for _, value := range mySlice {
		m[&value] = value // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}

func rangeValAddress6() {
	type v struct {
		id string
	}
	m := []*string{}

	mySlice := []v{{id: "A"}, {id: "B"}, {id: "C"}}
	for _, value := range mySlice {
		m = append(m, &value.id) // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}

func rangeValAddress7() {
	type v struct {
		id string
	}
	m := []*string{}

	for _, value := range []v{{id: "A"}, {id: "B"}, {id: "C"}} {
		m = append(m, &value.id) // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}

func rangeValAddress8() {
	type v struct {
		id string
	}
	m := []*string{}

	mySlice := []*v{{id: "A"}, {id: "B"}, {id: "C"}}
	for _, value := range mySlice {
		m = append(m, &value.id)
	}
}

func rangeValAddress9() {
	type v struct {
		id string
	}
	m := []*string{}

	mySlice := map[string]*v{"a": {id: "A"}, "b": {id: "B"}, "c": {id: "C"}}
	for _, value := range mySlice {
		m = append(m, &value.id)
	}
}

func rangeValAddress10() {
	type v struct {
		id string
	}
	m := []*string{}

	for _, value := range map[string]*v{"a": {id: "A"}, "b": {id: "B"}, "c": {id: "C"}} {
		m = append(m, &value.id)
	}
}

func rangeValAddress11() {
	type v struct {
		id string
	}
	m := map[string]*string{}

	for key, value := range map[string]*v{"a": {id: "A"}, "b": {id: "B"}, "c": {id: "C"}} {
		m[key] = &value.id
	}
}

func rangeValAddress12() {
	type v struct {
		id string
	}
	m := map[string]*string{}

	for key, value := range map[string]v{"a": {id: "A"}, "b": {id: "B"}, "c": {id: "C"}} {
		m[key] = &value.id // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}

func rangeValAddress13() {
	type v struct {
		id string
	}
	m := []*string{}

	otherSlice := map[string]*v{"a": {id: "A"}, "b": {id: "B"}, "c": {id: "C"}}
	mySlice := otherSlice
	for _, value := range mySlice {
		m = append(m, &value.id)
	}
}

func rangeValAddress14() {
	type v struct {
		id *string
	}

	m := []v{}
	for _, value := range []string{"A", "B", "C"} {
		a := v{id: &value} // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
		m = append(m, a)
	}
}

func rangeValAddress15() {
	type v struct {
		id *string
	}

	m := []v{}
	for _, value := range []string{"A", "B", "C"} {
		m = append(m, v{id: &value}) // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}

func rangeValAddressNestedField() {
	type nested struct {
		field struct {
			value string
		}
	}

	var pointers []*string
	for _, value := range []nested{{}} {
		pointers = append(pointers, &value.field.value) // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}
