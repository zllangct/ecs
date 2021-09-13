package main

import (
	"ecs"
	"fmt"
	"strconv"
)

func main()  {
	type Item struct {
		Count int
		Name  string
	}
	caseCount := 100
	var srcList []Item
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, Item{
			Count: i,
			Name:  "foo" + strconv.Itoa(i),
		})
	}

	c := ecs.NewCollection[Item]()

	cmp := map[int64]int{}
	for i := 0; i < caseCount; i++ {
		id, _ := c.Add(&srcList[i])
		cmp[id] = i
	}

	ret := c.Get(5)
	println(ret.Name)

	for iter := ecs.NewIterator(c) ; !iter.End(); iter.Next(){
		v := iter.Val()
		fmt.Printf("%+v", v)
	}
}
