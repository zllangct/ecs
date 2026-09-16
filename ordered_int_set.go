package ecs

type OrderedIntSet[T Integer] []T

func (c *OrderedIntSet[T]) FindIndexToInsert(it T, offset int) int {
	if len(*c) == 0 {
		return 0
	}
	l := 0
	r := len(*c) - 1
	if offset < r {
		l = offset
	}
	m := 0
	for l < r {
		m = (l + r) / 2
		if (*c)[m] > it {
			r = m - 1
		} else if (*c)[m] < it {
			l = m + 1
		} else {
			return -1
		}
	}
	if (*c)[l] < it {
		l = l + 1
	} else if (*c)[l] > it {
	} else {
		// 重复元素：返回 -1，调用方不得插入
		return -1
	}
	return l
}

func (c *OrderedIntSet[T]) Find(it T) (int, bool) {
	if len(*c) == 0 || it < (*c)[0] || it > (*c)[len(*c)-1] {
		return 0, false
	}
	l := 0
	r := len(*c) - 1
	m := 0
	for l <= r {
		m = (l + r) / 2
		if (*c)[m] == it {
			return m, true
		} else if (*c)[m] > it {
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return 0, false
}

func (c *OrderedIntSet[T]) Exist(it T) bool {
	if len(*c) == 0 {
		return false
	}
	_, ok := c.Find(it)
	return ok
}

func (c *OrderedIntSet[T]) IsSubSet(subSet OrderedIntSet[T]) bool {
	offset := 0
	length := len(*c)
	exist := false
	var temp T
	for i := 0; i < len(subSet); i++ {
		exist = false
		temp = subSet[i]
		for j := offset; j < length; j++ {
			if (*c)[j] > temp {
				return false
			}
			if (*c)[j] == temp {
				offset = j + 1
				exist = true
				break
			}
		}
		if !exist {
			return false
		}
	}
	return true
}

func (c *OrderedIntSet[T]) insert(it T, idx int) {
	*c = append(*c, T(0))
	copy((*c)[idx+1:], (*c)[idx:len(*c)-1])
	(*c)[idx] = it
}

func (c *OrderedIntSet[T]) Add(it T) bool {
	idx := c.FindIndexToInsert(it, 0)
	if idx < 0 {
		return false
	}
	c.insert(it, idx)
	return true
}

func (c *OrderedIntSet[T]) Merge(otherSet OrderedIntSet[T]) {
	if len(*c) == 0 {
		*c = otherSet
		return
	}
	index := 0
	for i := 0; i < len(otherSet); i++ {
		idx := c.FindIndexToInsert(otherSet[i], index)
		if idx < 0 {
			continue // 已存在，跳过
		}
		c.insert(otherSet[i], idx)
		index = idx
	}
}

func (c *OrderedIntSet[T]) Remove(it T) bool {
	idx, ok := c.Find(it)
	if !ok {
		return false
	}
	*c = append((*c)[:idx], (*c)[idx+1:]...)
	return true
}
