package main

import "math"

// Distance2D 非严格计算，仅做示范
func Distance2D(p1 *Position, p2 *Position) int {
	return int(math.Sqrt(math.Pow(2, float64(p1.X - p2.X)) + math.Pow(2, float64(p1.Y - p2.Y))))
}
