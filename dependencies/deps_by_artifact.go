package dependencies

import (
	. "github.com/dmolesUC3/mrt-build-info/jenkins"
	. "github.com/dmolesUC3/mrt-build-info/maven"
	"sort"
)

type Dependencies interface {
	ArtifactGraph() (ArtifactGraph, []error)
}

func DependenciesForJobs(jobs []Job) Dependencies {
	return &dependencies{jobs: jobs}
}

type dependencies struct {
	jobs []Job

	jobsByPom     map[Pom]Job
	poms          []Pom
	artifactGraph ArtifactGraph
}

func (d *dependencies) Poms() ([]Pom, []error) {
	var errors []error
	if d.poms == nil {
		var poms []Pom
		jobsByPom, errs := d.JobsByPom()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		for pom := range jobsByPom {
			poms = append(poms, pom)
		}
		sort.Sort(PomsByLocation(poms))
		d.poms = poms
	}
	return d.poms, errors
}

func (d *dependencies) JobsByPom() (map[Pom]Job, []error) {
	var errors []error
	if d.jobsByPom == nil {
		jobsByPom, errs := mapPomsToJobs(d.jobs)
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		d.jobsByPom = jobsByPom
	}
	return d.jobsByPom, errors
}

func (d *dependencies) ArtifactGraph() (ArtifactGraph, []error) {
	var errors []error
	if d.artifactGraph == nil {
		poms, errs := d.Poms()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		graph, errs := NewGraph(poms)
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		d.artifactGraph = graph
	}
	return d.artifactGraph, errors
}
