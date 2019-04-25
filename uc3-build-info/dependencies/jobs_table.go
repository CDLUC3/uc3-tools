package dependencies

import (
	. "github.com/dmolesUC3/mrt-build-info/jenkins"
	"github.com/dmolesUC3/mrt-build-info/shared"
)

func JobsTable(g JobGraph) shared.Table {
	model := &jobsTableModel{graph: g}
	return newTable(model)
}

type jobsTableModel struct {
	graph JobGraph
}

func (m *jobsTableModel) Rows() int {
	return len(m.graph.Jobs())
}

func (m *jobsTableModel) ItemType() string {
	return "Job"
}

func (m *jobsTableModel) ItemAt(row int) string {
	return m.jobAt(row).Name()
}

func (m *jobsTableModel) DependenciesOf(row int) string {
	job := m.jobAt(row)
	deps, _ := m.graph.DependenciesOf(job)
	// TODO: log errors
	return JobsByName(deps).String()
}

func (m *jobsTableModel) DependenciesOn(row int) string {
	job := m.jobAt(row)
	deps, _ := m.graph.DependenciesOn(job)
	// TODO: log errors
	return JobsByName(deps).String()
}

// ------------------------------------------------------------
// Unexported symbols

func (m *jobsTableModel) jobAt(row int) Job {
	return m.graph.Jobs()[row]
}


