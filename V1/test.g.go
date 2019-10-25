package main

type Test struct {
	P *map[int]int
	mP **map[int]int
	m1 map[int]int
	m2 map[int]int
}

func main()  {
	s:=[]int{1,2,3}
	s = s[:len(s)-1]
	println(s)
}
