package main

import "fmt"

const myConst = 42

const (
	iA = iota
	iB = iota
	iC = iota
)

func main() {
	var a int = 10
	var b string = "text"
	var c bool = true

	d := 3.14
	e := "auto type"
	f := false

	var (
		x int    = 111
		y string = "text"
		z byte   = 255
	)

	fmt.Println(a, b, c, d, e, f, x, y, z, iA, iB, iC)
}
