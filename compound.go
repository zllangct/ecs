package ecs

import "errors"

type Compound []uint16

func (c Compound) insertIndex(it uint16) int {
	l := 0
	r := len(c) - 1
	m := 0
	for l < r {
		m = (l + r) / 2
		if c[m] > it {
			r = m - 1
		} else if c[m] < it {
			l = m + 1
		} else {
			return -1
		}
	}
	if c[l] < it {
		l = l + 1
	} else if c[l] > it {
	} else {
		l = l - 1
	}
	return l
}

func (c Compound) find(it uint16) int {
	l := 0
	r := len(c) - 1
	m := 0
	for l <= r {
		m = (l + r) / 2
		if c[m] == it {
			return m
		} else if c[m] > it {
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return -1
}

func (c *Compound) Add(it uint16) error {
	idx := c.insertIndex(it)
	if idx < 0 {
		return errors.New("this type already existed")
	}
	*c = append(*c, 0)
	copy((*c)[idx+1:], (*c)[idx:len(*c)-1])
	(*c)[idx] = it
	return nil
}

func (c *Compound) Remove(it uint16) {
	idx := c.find(it)
	if idx < 0 {
		return
	}
	*c = append((*c)[:idx], (*c)[idx+1:]...)
}

func (c Compound) Type() interface{} {
	if len(c) == 0 {
		return nil
	}
	return getCompoundType(c)
}
