package ecs

type OrderedIntSet[T Integer] []T

func (c *OrderedIntSet[T]) InsertIndex(it T) int {
	if len(*c) == 0 {
		return 0
	}
	l := 0
	r := len(*c) - 1
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
		l = l - 1
	}
	return l
}

func (c *OrderedIntSet[T]) Find(it T) int {
	l := 0
	r := len(*c) - 1
	m := 0
	for l <= r {
		m = (l + r) / 2
		if (*c)[m] == it {
			return m
		} else if (*c)[m] > it {
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return -1
}

func (c *OrderedIntSet[T]) Exist(it T) bool {
	return c.Find(it) != -1
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

func (c *OrderedIntSet[T]) Add(it T) bool {
	idx := c.InsertIndex(it)
	if idx < 0 {
		return false
	}
	*c = append(*c, 0)
	copy((*c)[idx+1:], (*c)[idx:len(*c)-1])
	(*c)[idx] = it
	return true
}

func (c *OrderedIntSet[T]) Remove(it T) bool {
	idx := c.Find(it)
	if idx < 0 {
		return false
	}
	*c = append((*c)[:idx], (*c)[idx+1:]...)
	return true
}
