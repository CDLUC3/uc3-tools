package shared

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	if Flags.Verbose && !Flags.TSV {
		fmt.Fprintf(os.Stderr, "Formatting %d rows", t.Rows())
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

		if Flags.Verbose && !Flags.TSV {
			fmt.Fprint(os.Stderr, ".")
		}
	}
	if Flags.Verbose && !Flags.TSV {
		fmt.Fprintln(os.Stderr, "Done.")
	}
}
