package main

import "fmt"

type Test struct {
	Foo string `json:"foo"`
	bar string `json:"bar"`
}

func main() {
	t := Test{"baz", "quux"}
	fmt.Println(t)
	fmt.Println("Hello Aeneas")
}
