package main

import (
	"fmt"
	"time"
)

func main() {
	t1 := time.Now()
	t2 := time.Now()
	fmt.Println(t1 == t2)
}
