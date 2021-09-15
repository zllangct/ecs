package main

import "math"

type Point struct {
	X int
	Y int
	Z int
}

func Distance(p1 Point, p2 Point) int {
	return int(math.Pow(2, float64(p1.X - p2.X)))
}
