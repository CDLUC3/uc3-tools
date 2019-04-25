package dependencies

import (
	. "github.com/CDLUC3/uc3-tools/mrt-build-info/jenkins"
	. "github.com/CDLUC3/uc3-tools/mrt-build-info/maven"
	"github.com/CDLUC3/uc3-tools/mrt-build-info/shared"
)

func PomsTable(g JobGraph) (shared.Table, []error) {
	pgraph, errs := g.PomGraph()
	if pgraph == nil {
		return nil, errs
	}
	model := &pomsTableModel{graph: pgraph}
	return newTable(model), errs
}

type pomsTableModel struct {
	graph PomGraph
}

func (m *pomsTableModel) Rows() int {
	return len(m.graph.Poms())
}

func (m *pomsTableModel) ItemType() string {
	return "Pom"
}

func (m *pomsTableModel) ItemAt(row int) string {
	return m.pomAt(row).String()
}

func (m *pomsTableModel) DependenciesOf(row int) string {
	pom := m.pomAt(row)
	deps := m.graph.DependenciesOf(pom)
	return PomsByLocation(deps).String()
}

func (m *pomsTableModel) DependenciesOn(row int) string {
	pom := m.pomAt(row)
	deps := m.graph.DependenciesOn(pom)
	return PomsByLocation(deps).String()
}

// ------------------------------------------------------------
// Unexported symbols

func (m *pomsTableModel) pomAt(row int) Pom {
	return m.graph.Poms()[row]
}


