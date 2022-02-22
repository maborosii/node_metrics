package main

import "fmt"

func main() {

	slice := []int{0, 1, 2, 3}
	fmt.Printf("slice: %v slice addr %p \n", slice, &slice)

	ret := changeSlice(slice)
	fmt.Printf("slice: %v ret: %v slice addr %p \n", slice, ret, &slice)
}

func changeSlice(s []int) []int {
	s[1] = 111
	fmt.Printf("%p\n", &s)
	return s
}
