package dependencies

import (
	. "github.com/CDLUC3/uc3-tools/uc3-build-info/jenkins"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/shared"
)

func JobsTable(g JobGraph, showCounts bool) shared.Table {
	model := &jobsTableModel{
		graph:      g,
		showCounts: showCounts,
	}
	return newTable(model)
}

type jobsTableModel struct {
	graph      JobGraph
	showCounts bool
}

func (m *jobsTableModel) Rows() int {
	return len(m.graph.Jobs())
}

func (m *jobsTableModel) ItemType() string {
	return "Job"
}

func (m *jobsTableModel) ItemAt(row int) string {
	return m.JobName(row)
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

func (m *jobsTableModel) CountDependenciesOf(row int) int {
	job := m.jobAt(row)
	deps, _ := m.graph.DependenciesOf(job)
	// TODO: log errors
	return len(deps)
}

func (m *jobsTableModel) CountDependenciesOn(row int) int {
	job := m.jobAt(row)
	deps, _ := m.graph.DependenciesOn(job)
	// TODO: log errors
	return len(deps)
}

func (m *jobsTableModel) ShowCounts() bool {
	return m.showCounts
}

func (m *jobsTableModel) ShowJobs() bool {
	return false
}

func (m *jobsTableModel) JobName(row int) string {
	return m.jobAt(row).Name()
}

// ------------------------------------------------------------
// Unexported symbols

func (m *jobsTableModel) jobAt(row int) Job {
	return m.graph.Jobs()[row]
}
