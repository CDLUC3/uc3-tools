package maven

import (
	"fmt"
	"sort"
)

// ------------------------------------------------------------
// Graph

type Graph interface {
	Artifacts() []Artifact
	SortedArtifacts() []Artifact
	DependenciesOf(artifact Artifact) (deps []Artifact)
	DependenciesOn(artifact Artifact) (deps []Artifact)
	PomForArtifact(artifact Artifact) Pom
}

func NewGraph(poms []Pom) (Graph, []error) {
	var errors []error
	artifacts, pomsByArtifact, errs1 := artifactsFromPoms(poms)
	errors = append(errors, errs1...)

	depsByFrom, depsByTo, errs2 := dependencies(poms, pomsByArtifact)
	errors = append(errors, errs2...)

	g := graph{artifacts: artifacts, pomsByArtifact: pomsByArtifact, depsByFrom: depsByFrom, depsByTo: depsByTo}
	return &g, errors
}

// ------------------------------------------------------------
// Unexported symbols

type dependency struct {
	fromArtifact Artifact
	toArtifact   Artifact
}

type graph struct {
	artifacts      []Artifact
	pomsByArtifact map[Artifact]Pom
	depsByFrom     map[Artifact][]dependency
	depsByTo       map[Artifact][]dependency

	sortedArtifacts []Artifact
}

func (g *graph) DependenciesOf(artifact Artifact) (deps []Artifact) {
	if depsFromArtifact, ok := g.depsByFrom[artifact]; ok {
		for _, d := range depsFromArtifact {
			deps = append(deps, d.toArtifact)
		}
	}
	return deps
}

func (g *graph) DependenciesOn(artifact Artifact) (deps []Artifact) {
	if depsToArtifact, ok := g.depsByTo[artifact]; ok {
		for _, d := range depsToArtifact {
			deps = append(deps, d.toArtifact)
		}
	}
	return deps
}

func (g *graph) Artifacts() []Artifact {
	return g.artifacts
}

func (g *graph) SortedArtifacts() []Artifact {
	if g.sortedArtifacts == nil {
		g.sortedArtifacts = newTopoSort(g).sortedArtifacts()
	}
	return g.sortedArtifacts
}

func (g *graph) PomForArtifact(artifact Artifact) Pom {
	return g.pomsByArtifact[artifact]
}

func artifactsFromPoms(poms []Pom) ([]Artifact, map[Artifact]Pom, []error) {
	var errors []error

	var allArtifacts []Artifact
	a2p := map[Artifact]Pom{}
	for _, pom := range poms {
		artifact, err := pom.Artifact()
		if err != nil {
			errors = append(errors, err)
			continue
		}
		allArtifacts = append(allArtifacts, artifact)
		a2p[artifact] = pom
	}
	sort.Sort(ArtifactsByString(allArtifacts))
	return allArtifacts, a2p, errors
}

func dependencies(allPoms []Pom, pomsByArtifact map[Artifact]Pom) (depsByFrom map[Artifact][]dependency, depsByTo map[Artifact][]dependency, err []error) {
	var errors []error

	depsByFrom = map[Artifact][]dependency{}
	depsByTo = map[Artifact][]dependency{}
	for _, pom := range allPoms {
		fromArtifact, err := pom.Artifact()
		if err != nil {
			errors = append(errors, err)
			continue
		}
		pomDeps, depErrs := pom.Dependencies()
		errors = append(errors, depErrs...)

		for _, toArtifact := range pomDeps {
			var ok bool
			if _, ok = pomsByArtifact[toArtifact]; !ok {
				// third-party dependency
				continue
			}
			var fromDeps []dependency
			if fromDeps, ok = depsByFrom[fromArtifact]; !ok {
				fromDeps = []dependency{}
			}
			var toDeps []dependency
			if toDeps, ok = depsByTo[toArtifact]; !ok {
				toDeps = []dependency{}
			}
			dep := dependency{fromArtifact: fromArtifact, toArtifact: toArtifact}
			depsByFrom[fromArtifact] = append(fromDeps, dep)
			depsByTo[toArtifact] = append(toDeps, dep)
		}
	}
	return depsByFrom, depsByTo, errors
}

// ------------------------------------------------------------
// Topological sorter

type topoSort struct {
	g        *graph
	sorted   []Artifact
	visiting map[Artifact]bool
	visited  map[Artifact]bool
}

func newTopoSort(g *graph) *topoSort {
	return &topoSort{
		g:        g,
		sorted:   []Artifact{},
		visiting: map[Artifact]bool{},
		visited:  map[Artifact]bool{},
	}
}

func (t *topoSort) markVisiting(n Artifact) {
	if t.visiting == nil {
		t.visiting = map[Artifact]bool{}
	}
	t.visiting[n] = true
}

func (t *topoSort) unmarkVisiting(n Artifact) {
	if t.visiting == nil {
		panic(fmt.Errorf("can't unmark artifact %v (never marked)", n))
	}
	t.visiting[n] = false
}

func (t *topoSort) markVisited(n Artifact) {
	if t.visited == nil {
		t.visited = map[Artifact]bool{}
	}
	t.visited[n] = true
}

func (t *topoSort) visit(n Artifact) {
	if t.visited[n] {
		return
	}
	if t.visiting[n] {
		panic("not a DAG")
	}
	t.markVisiting(n)
	for _, m := range t.g.DependenciesOf(n) {
		t.visit(m)
	}
	t.unmarkVisiting(n)
	t.markVisited(n)
	t.sorted = append(t.sorted, n)
}

func (t *topoSort) sortedArtifacts() []Artifact {
	for _, n := range t.g.artifacts {
		if t.visited[n] {
			continue
		}
		t.visit(n)
	}
	for l, r := 0, len(t.sorted)-1; l < r; l, r = l+1, r-1 {
		t.sorted[l], t.sorted[r] = t.sorted[r], t.sorted[l]
	}
	return t.sorted
}
