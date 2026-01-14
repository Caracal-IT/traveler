package main

import (
	"fmt"
	colors "traveler/play_ground/basic/colours"
)

type Color int

const (
	Red Color = iota
	Green
	Blue
)

func main() {
	const (
		a = iota * 2
		b
		c
	)

	fmt.Println(a, b, c)

	var red = Red
	fmt.Println(red)

	fmt.Println(colors.Green)
}
