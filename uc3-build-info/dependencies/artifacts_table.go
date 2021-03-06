package dependencies

import (
	. "github.com/CDLUC3/uc3-tools/uc3-build-info/jenkins"
	. "github.com/CDLUC3/uc3-tools/uc3-build-info/maven"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/shared"
)

func ArtifactsTable(g JobGraph, showCounts bool, showJobs bool) (shared.Table, []error) {
	agraph, errs := g.ArtifactGraph()
	if agraph == nil {
		return nil, errs
	}
	model := &artifactsTableModel{
		jobGraph:   g,
		agraph:     agraph,
		showCounts: showCounts,
		showJobs: showJobs,
	}
	return newTable(model), errs
}

type artifactsTableModel struct {
	jobGraph   JobGraph
	agraph     ArtifactGraph
	showCounts bool
	showJobs bool
}

func (m *artifactsTableModel) Rows() int {
	return len(m.artifacts())
}

func (m *artifactsTableModel) ItemType() string {
	return "Artifact"
}

func (m *artifactsTableModel) ItemAt(row int) string {
	return m.artifactAt(row).String()
}

func (m *artifactsTableModel) DependenciesOf(row int) string {
	artifact := m.artifactAt(row)
	deps := m.agraph.DependenciesOf(artifact)
	return ArtifactsByString(deps).String()
}

func (m *artifactsTableModel) DependenciesOn(row int) string {
	artifact := m.artifactAt(row)
	deps := m.agraph.DependenciesOn(artifact)
	return ArtifactsByString(deps).String()
}

func (m *artifactsTableModel) CountDependenciesOf(row int) int {
	artifact := m.artifactAt(row)
	deps := m.agraph.DependenciesOf(artifact)
	return len(deps)
}

func (m *artifactsTableModel) CountDependenciesOn(row int) int {
	artifact := m.artifactAt(row)
	deps := m.agraph.DependenciesOn(artifact)
	return len(deps)
}

func (m *artifactsTableModel) ShowCounts() bool {
	return m.showCounts
}

func (m *artifactsTableModel) ShowJobs() bool {
	return m.showJobs
}

func (m *artifactsTableModel) JobName(row int) string {
	artifact := m.artifactAt(row)
	pom := m.agraph.PomForArtifact(artifact)
	job, _ := m.jobGraph.JobFor(pom)
	// TODO: log errors
	if job == nil {
		return shared.ValueUnknown
	}
	return job.Name()
}

// ------------------------------------------------------------
// Unexported symbols

func (m *artifactsTableModel) artifacts() []Artifact {
	return m.agraph.SortedArtifacts()
}

func (m *artifactsTableModel) artifactAt(row int) Artifact {
	return m.artifacts()[row]
}
