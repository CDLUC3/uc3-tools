package shared

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type Table interface {
	Cols() int
	Rows() int
	HeaderFor(col int) string
	ValueAt(row, col int) string
	Print(w io.Writer, sep string)
}

func TableFrom(cols ...TableColumn) Table {
	return &table{columns: cols}
}

// ------------------------------------------------------------
// Unexported symbols

type table struct {
	columns []TableColumn
	rows    *int
}

func (t *table) Cols() int {
	return len(t.columns)
}

func (t *table) Rows() int {
	if t.rows == nil {
		rows := 0
		for _, col := range t.columns {
			colRows := col.Rows()
			if colRows > rows {
				rows = colRows
			}
		}
		t.rows = &rows
	}
	return *t.rows
}

func (t *table) HeaderFor(col int) string {
	return t.columns[col].Header()
}

func (t *table) ValueAt(row, col int) string {
	column := t.columns[col]
	if row >= column.Rows() {
		return ""
	}
	return column.ValueAt(row)
}

//noinspection GoUnhandledErrorResult
func (t *table) Print(w io.Writer, sep string) {
	if t == nil {
		fmt.Fprintln(w, "(nil table)")
	}

	var out *bufio.Writer
	if Flags.TSV {
		out = bufio.NewWriter(w)
	} else {
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
		defer tw.Flush()
		out = bufio.NewWriter(tw)
	}

	cols := t.Cols()
	for col := 0; col < cols; col++ {
		out.WriteString(t.HeaderFor(col))
		if col+1 < cols {
			out.WriteString(sep)
		}
	}
	out.WriteRune('\n')
	out.Flush()

	rows := t.Rows()
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			value := t.ValueAt(row, col)
			out.WriteString(value)
			if col+1 < cols {
				out.WriteString(sep)
			}
		}
		out.WriteRune('\n')
		out.Flush()
	}
}

func SplitRows(t Table, sep string) Table {
	var splitRows [][]string

	cols := t.Cols()
	rows := t.Rows()
	for row := 0; row <= rows; row++ {
		cells := make([]string, cols)
		for col := 0; col < cols; col++ {
			cells[col] = t.ValueAt(row, col)
		}
		splitRows = append(splitRows, splitRow(cells, sep)...)
	}

	splitRowCount := len(splitRows)
	for i, srow := range splitRows {
		if len(srow) != cols {
			panic(fmt.Errorf("%d: expected %d cells, found %d: %#v", i, cols, len(srow), srow))
		}
	}

	columns := make([]TableColumn, cols)
	for c := 0; c < cols; c++ {
		cIndex := c // make sure we close over the current value
		columns[cIndex] = NewTableColumn(t.HeaderFor(cIndex), splitRowCount, func(row int) string {
			srow := splitRows[row]
			return srow[cIndex]
		})
	}
	return &table{columns: columns, rows: &splitRowCount}
}

/*
splitRow splits a row vertically based on a separator, e.g.

 	{"a", "b", "c,d,e", "f,g"}

becomes

	{
		{"a", "b", "c", "f"},
		{"", "", "d", "g"},
		{"", "", "e", ""},
	}
*/
func splitRow(cells []string, sep string) [][]string {
	var columns [][]string
	var rowCount = 0
	for _, cell := range cells {
		cellCol := strings.Split(cell, sep)
		columns = append(columns, cellCol)
		if len(cellCol) > rowCount {
			rowCount = len(cellCol)
		}
	}

	cols := len(cells)
	rows := make([][]string, rowCount)
	for row := 0; row < rowCount; row++ {
		sRow := make([]string, cols)
		for col := 0; col < cols; col++ {
			var cellCol = columns[col]
			if len(cellCol) > row {
				cellVal := cellCol[row]
				sRow[col] = strings.TrimSpace(cellVal)
			}
		}
		rows[row] = sRow
	}

	return rows
}


