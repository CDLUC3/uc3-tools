package dependencies

import (
	. "github.com/CDLUC3/uc3-tools/uc3-build-info/jenkins"
	. "github.com/CDLUC3/uc3-tools/uc3-build-info/maven"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/shared"
)

func PomsTable(g JobGraph, showCounts bool) (shared.Table, []error) {
	pgraph, errs := g.PomGraph()
	if pgraph == nil {
		return nil, errs
	}
	model := &pomsTableModel{
		graph:      pgraph,
		showCounts: showCounts,
	}
	return newTable(model), errs
}

type pomsTableModel struct {
	graph      PomGraph
	showCounts bool
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

func (m *pomsTableModel) CountDependenciesOf(row int) int {
	pom := m.pomAt(row)
	deps := m.graph.DependenciesOf(pom)
	return len(deps)
}

func (m *pomsTableModel) CountDependenciesOn(row int) int {
	pom := m.pomAt(row)
	deps := m.graph.DependenciesOn(pom)
	return len(deps)
}

func (m *pomsTableModel) ShowCounts() bool {
	return m.showCounts
}

// ------------------------------------------------------------
// Unexported symbols

func (m *pomsTableModel) pomAt(row int) Pom {
	return m.graph.Poms()[row]
}
