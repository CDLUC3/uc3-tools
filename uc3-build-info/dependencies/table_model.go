package dependencies

import (
	. "github.com/dmolesUC3/mrt-build-info/shared"
)

type tableModel interface {
	Rows() int
	ItemType() string
	ItemAt(row int) string
	DependenciesOf(row int) string
	DependenciesOn(row int) string
}

func newTable(m tableModel) Table {
	rows := m.Rows()
	return TableFrom(
		NewTableColumn(m.ItemType(), rows, m.ItemAt),
		NewTableColumn("Requires", rows, m.DependenciesOf),
		NewTableColumn("Required by", rows, m.DependenciesOn),
	)
}

