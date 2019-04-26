package shared

import (
	"fmt"
	"strconv"
)

type TableColumn interface {
	Header() string
	Rows() int
	ValueAt(row int) string
}

func NewTableColumn(header string, rows int, valueAt func(row int) string) TableColumn {
	return &tableColumn{header: header, rows: rows, valueAt: valueAt}
}

func NewIntTableColumn(header string, rows int, valueAt func(row int) int) TableColumn {
	stringValueAt := func(row int) string {
		return strconv.Itoa(valueAt(row))
	}
	return NewTableColumn(header, rows, stringValueAt)
}

// ------------------------------------------------------------
// Unexported symbols

type tableColumn struct {
	header string
	rows int
	valueAt func(row int) string
}

func (c *tableColumn) Header() string {
	return c.header
}

func (c *tableColumn) Rows() int {
	return c.rows
}

func (c *tableColumn) ValueAt(row int) string {
	if row >= c.rows {
		panic(fmt.Errorf("row out of bounds: %d of %d", row, c.rows))
	}
	return c.valueAt(row)
}



