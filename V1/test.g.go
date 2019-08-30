package main

type testList []int

func (p *testList)Insert(system int,index int)  {
	if len(*p)==0 {
		*p = append(*p, system)
	}
	for i, _ := range *p {
		if index == i {
			break
		}else if index < i {
			*p = append(append((*p)[:i-1], system), (*p)[i:]...)
		}
	}

}
func main()  {
	i:=[]int{1,2,3,4,5,6}
	s:=append([]int{},i[2:]...)
	i=append(append(i[:2], 0), s...)
}
