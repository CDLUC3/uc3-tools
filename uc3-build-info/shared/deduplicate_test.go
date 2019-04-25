package shared

import (
	. "gopkg.in/check.v1"
	"sort"
	"strconv"
)

type DeduplicateSuite struct{}

var _ = Suite(&DeduplicateSuite{})

type IntsByValue []int

func (in IntsByValue) Len() int {
	return len(in)
}

func (in IntsByValue) Less(i, j int) bool {
	return in[i] < in[j]
}

func (in IntsByValue) Swap(i, j int) {
	in[i], in[j] = in[j], in[i]
}

func (in IntsByValue) Eq(i, j int) bool {
	return in[i] == in[j]
}

func (in IntsByValue) Copy(j, i int) {
	in[j] = in[i]
}

func (s *DeduplicateSuite) TestDeduplicate(c *C) {
	in := []int{3, 2, 1, 4, 3, 2, 1, 4, 1}
	Deduplicate(IntsByValue(in), func(len int) { in = in[:len] })
	c.Assert(len(in), Equals, 4)
	for i := 0; i < 4; i++ {
		c.Check(in[i], Equals, 1+i)
	}
}

func (s *DeduplicateSuite) TestDeduplicateAny(c *C) {
	in := []string{"3", "2", "1", "4", "3", "2", "1", "4", "1"}
	DeduplicateAny(
		func() { sort.Strings(in) },
		func() int { return len(in) },
		func(i, j int) bool { return in[i] == in[j] },
		func(j, i int) { in[j] = in[i] },
		func(len int) { in = in[:len] },
	)
	c.Assert(len(in), Equals, 4)
	for i := 0; i < 4; i++ {
		c.Check(in[i], Equals, strconv.Itoa(1+i))
	}
}
