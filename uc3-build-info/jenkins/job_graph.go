package jenkins

import (
	"fmt"
	. "github.com/CDLUC3/uc3-tools/uc3-build-info/maven"
	. "github.com/CDLUC3/uc3-tools/uc3-build-info/shared"
	"os"
	"sort"
)

type JobGraph interface {
	ArtifactGraph() (ArtifactGraph, []error)
	PomGraph() (PomGraph, []error)
	Jobs() []Job
	JobFor(pom Pom) (Job, []error)
	DependenciesOf(job Job) (deps []Job, errors []error)
	DependenciesOn(job Job) (deps []Job, errors []error)
}

func NewJobGraph(jobs []Job) JobGraph {
	return &jobGraph{jobs: jobs}
}

type jobDep struct {
	fromJob Job
	toJob   Job
}

type jobGraph struct {
	jobs []Job

	jobsByPom map[Pom]Job
	poms      []Pom

	artifactGraph ArtifactGraph
	pomGraph      PomGraph

	allDeps    []jobDep
	depsByFrom map[Job][]jobDep
	depsByTo   map[Job][]jobDep
}

// TODO: SortedJobs()?
func (g *jobGraph) Jobs() []Job {
	return g.jobs
}

func (g *jobGraph) Poms() ([]Pom, []error) {
	var errors []error
	if g.poms == nil {
		var poms []Pom
		jobsByPom, errs := g.JobsByPom()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		for pom := range jobsByPom {
			poms = append(poms, pom)
		}
		sort.Sort(PomsByLocation(poms))
		g.poms = poms
	}
	return g.poms, errors
}

func (g *jobGraph) JobFor(pom Pom) (Job, []error) {
	jobsByPom, errs := g.JobsByPom()
	if jobsByPom != nil {
		if job, ok := jobsByPom[pom]; ok {
			return job, errs
		}
	}
	errs = append(errs, fmt.Errorf("no job found for pom %v", pom))
	return nil, errs
}

func (g *jobGraph) JobsByPom() (map[Pom]Job, []error) {
	var errors []error
	if g.jobsByPom == nil {
		jobsByPom, errs := mapPomsToJobs(g.jobs)
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		g.jobsByPom = jobsByPom
	}
	return g.jobsByPom, errors
}

func (g *jobGraph) ArtifactGraph() (ArtifactGraph, []error) {
	var errors []error
	if g.artifactGraph == nil {
		poms, errs := g.Poms()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		graph, errs := NewArtifactGraph(poms)
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		g.artifactGraph = graph
	}
	return g.artifactGraph, errors
}

func (g *jobGraph) PomGraph() (PomGraph, []error) {
	var errors []error
	if g.pomGraph == nil {
		agraph, errs := g.ArtifactGraph()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		g.pomGraph = NewPomGraph(agraph)
	}
	return g.pomGraph, errors
}

func (g *jobGraph) DependenciesOf(job Job) (deps []Job, errors []error) {
	_, depsByFrom, _, errs := g.deps()
	if depsByFrom != nil {
		if depsFromJob, ok := depsByFrom[job]; ok {
			for _, d := range depsFromJob {
				deps = append(deps, d.toJob)
			}
		}
	}
	Deduplicate(JobsByName(deps), func(len int) { deps = deps[:len] })
	return deps, errs
}

func (g *jobGraph) DependenciesOn(job Job) (deps []Job, errors []error) {
	_, _, depsByTo, errs := g.deps()
	if depsByTo != nil {
		if depsToJob, ok := depsByTo[job]; ok {
			for _, d := range depsToJob {
				deps = append(deps, d.fromJob)
			}
		}
	}
	Deduplicate(JobsByName(deps), func(len int) { deps = deps[:len] })
	return deps, errs
}

func (g *jobGraph) deps() (allDeps []jobDep, depsByFrom map[Job][]jobDep, depsByTo map[Job][]jobDep, errors []error) {
	if g.allDeps == nil {
		allDeps = []jobDep{}
		depsByFrom = map[Job][]jobDep{}
		depsByTo = map[Job][]jobDep{}

		pgraph, errs := g.PomGraph()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		if pgraph == nil {
			return
		}

		jobsByPom, errs := g.JobsByPom()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		if jobsByPom == nil {
			return
		}

		poms := pgraph.Poms()
		//noinspection GoUnhandledErrorResult
		if Flags.Verbose {
			fmt.Fprintf(os.Stderr, "Determining job dependencies for %d poms...", len(poms))
			defer func() { fmt.Fprintln(os.Stderr, "Done.") }()
		}
		for _, p := range poms {
			//noinspection GoUnhandledErrorResult
			if Flags.Verbose {
				fmt.Fprint(os.Stderr, ".")
			}
			fromPom := p
			fromJob := jobsByPom[fromPom]
			for _, toPom := range pgraph.DependenciesOf(fromPom) {
				toJob := jobsByPom[toPom]
				if fromJob == toJob {
					// Can happen, even w/o POM self-deps, for jobs that build multiple related POMs
					continue
				}

				var ok bool
				var fromDeps []jobDep
				if fromDeps, ok = depsByFrom[fromJob]; !ok {
					fromDeps = []jobDep{}
				}
				var toDeps []jobDep
				if toDeps, ok = depsByTo[toJob]; !ok {
					toDeps = []jobDep{}
				}

				dep := jobDep{fromJob: fromJob, toJob: toJob}
				allDeps = append(allDeps, dep)
				depsByFrom[fromJob] = append(fromDeps, dep)
				depsByTo[toJob] = append(toDeps, dep)
			}
		}
		g.allDeps = allDeps
		g.depsByFrom = depsByFrom
		g.depsByTo = depsByTo
	}
	return g.allDeps, g.depsByFrom, g.depsByTo, errors
}
