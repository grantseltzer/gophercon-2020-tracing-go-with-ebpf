package main

import "fmt"

//go:noinline
func addTwoNumbers(x int, y int8) {
	fmt.Printf("%d + %d = %d\n", x, y, x+int(y))
}

func main() {
	addTwoNumbers(42, 3)
}
