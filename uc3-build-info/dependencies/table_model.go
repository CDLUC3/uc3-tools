package dependencies

import (
	. "github.com/CDLUC3/uc3-tools/uc3-build-info/shared"
)

type tableModel interface {
	Rows() int
	ItemType() string
	ItemAt(row int) string
	DependenciesOf(row int) string
	DependenciesOn(row int) string
	CountDependenciesOf(row int) int
	CountDependenciesOn(row int) int
	ShowCounts() bool
}

func newTable(m tableModel) Table {
	rows := m.Rows()

	itemCol := NewTableColumn(m.ItemType(), rows, m.ItemAt)
	reqsCol := NewTableColumn("Requires", rows, m.DependenciesOf)
	reqdByCol := NewTableColumn("Required By", rows, m.DependenciesOn)

	cols := []TableColumn{ itemCol }
	if m.ShowCounts() {
		cols = append(cols, NewIntTableColumn("# Req'd By", rows, m.CountDependenciesOn))
	}
	cols = append(cols, reqdByCol)
	if m.ShowCounts() {
		cols = append(cols, NewIntTableColumn("# Reqs", rows, m.CountDependenciesOf))
	}
	cols = append(cols, reqsCol)

	return TableFrom(cols...)
}

