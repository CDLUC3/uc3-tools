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
	ShowJobs() bool
	JobName(row int) string
}

func newTable(m tableModel) Table {
	rows := m.Rows()

	cols := []TableColumn{NewTableColumn(m.ItemType(), rows, m.ItemAt)}
	if m.ShowJobs() {
		cols = append(cols, NewTableColumn("Job", rows, m.JobName))
	}
	if m.ShowCounts() {
		cols = append(cols, NewIntTableColumn("# Req'd By", rows, m.CountDependenciesOn))
	}
	cols = append(cols, NewTableColumn("Required By", rows, m.DependenciesOn))
	if m.ShowCounts() {
		cols = append(cols, NewIntTableColumn("# Reqs", rows, m.CountDependenciesOf))
	}
	cols = append(cols, NewTableColumn("Requires", rows, m.DependenciesOf))

	return TableFrom(cols...)
}
