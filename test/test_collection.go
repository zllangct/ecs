package main

import (
	"fmt"
	"github.com/zllangct/ecs"
	"strconv"
)

func main()  {
	//待存储的数据定义
	type Item struct {
		Count int
		Name  string
		Arr []int
	}

	//准备数据
	caseCount := 50
	var srcList []Item
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, Item{
			Count: i,
			Name:  "foo" + strconv.Itoa(i),
			Arr: []int{1,2,3},
		})
	}

	//创建容器(无序数据集)
	c := ecs.NewCollection[Item]()

	//添加数据
	cmp := map[int64]int{}
	for i := 0; i < caseCount; i++ {
		id, _ := c.Add(&srcList[i])
		cmp[id] = i
	}

	//遍历风格 1：
	for iter := ecs.NewIterator(c) ; !iter.End(); iter.Next(){
		v := iter.Val()
		fmt.Printf("style 1: %+v\n", v)
	}

	//遍历风格 2:
	iter := ecs.NewIterator(c)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		fmt.Printf("style 2: %+v\n", c)
	}
}
