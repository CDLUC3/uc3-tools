package shared

import (
	. "gopkg.in/check.v1"
)

type TableSuite struct{}

var _ = Suite(&TableSuite{})

func (s *TableSuite) TestSplitRows(c *C) {
	headers := []string{"1", "2", "3"}
	colCount := len(headers)
	data := [][]string{
		{"a", "b,c", "d"},
		{"e", "f,g", "h"},
	}
	columns := make([]TableColumn, colCount)
	for col := 0; col < colCount; col++ {
		colIndex := col
		columns[colIndex] = NewTableColumn(headers[col], len(data), func(row int) string {
			return data[row][colIndex]
		})
	}
	table := TableFrom(columns...)

	// just to be sure
	for row := 0; row < len(data); row++ {
		for col := 0; col < colCount; col++ {
			expected := data[row][col]
			actual := table.ValueAt(row, col)
			c.Check(actual, Equals, expected, Commentf("(%d, %d): expected %#v, was %#v", row, col, expected, actual))
		}
	}

	table = SplitRows(table, ",")
	data = [][]string{
		{"a", "b", "d"},
		{"", "c", ""},
		{"e", "f", "h"},
		{"", "g", ""},
	}

	for row := 0; row < len(data); row++ {
		for col := 0; col < colCount; col++ {
			expected := data[row][col]
			actual := table.ValueAt(row, col)
			c.Check(actual, Equals, expected, Commentf("(%d, %d): expected %#v, was %#v", row, col, expected, actual))
		}
	}
}
