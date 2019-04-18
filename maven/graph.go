package maven

import (
	"fmt"
	"sort"
)

// ------------------------------------------------------------
// Graph

type Graph interface {
	SortedArtifacts() []Artifact
	DependenciesOf(artifact Artifact) (deps []Artifact)
	DependenciesOn(artifact Artifact) (deps []Artifact)
	PomForArtifact(artifact Artifact) Pom
}

func NewGraph(poms []Pom) (Graph, error) {
	artifacts, pomsByArtifact, err := artifactsFromPoms(poms)
	if err != nil {
		return nil, err
	}

	depsByFrom, depsByTo, err := dependencies(poms, pomsByArtifact)
	if err != nil {
		return nil, err
	}

	g := graph{artifacts: artifacts, pomsByArtifact: pomsByArtifact, depsByFrom: depsByFrom, depsByTo: depsByTo}
	return &g, nil
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

func (g *graph) SortedArtifacts() []Artifact {
	if g.sortedArtifacts == nil {
		g.sortedArtifacts = newTopoSort(g).sortedArtifacts()
	}
	return g.sortedArtifacts
}

func (g *graph) PomForArtifact(artifact Artifact) Pom {
	return g.pomsByArtifact[artifact]
}

func artifactsFromPoms(poms []Pom) ([]Artifact, map[Artifact]Pom, error) {
	var allArtifacts []Artifact
	a2p := map[Artifact]Pom{}
	for _, pom := range poms {
		artifact, err := pom.Artifact()
		if err != nil {
			return nil, nil, err
		}
		allArtifacts = append(allArtifacts, artifact)
		a2p[artifact] = pom
	}
	sort.Sort(ArtifactsByString(allArtifacts))
	return allArtifacts, a2p, nil
}

func dependencies(allPoms []Pom, pomsByArtifact map[Artifact]Pom) (depsByFrom map[Artifact][]dependency, depsByTo map[Artifact][]dependency, err error) {
	depsByFrom = map[Artifact][]dependency{}
	depsByTo = map[Artifact][]dependency{}
	for _, pom := range allPoms {
		fromArtifact, err := pom.Artifact()
		if err != nil {
			return nil, nil, err // really shouldn't happen
		}
		pomDeps, err := pom.Dependencies()
		if err != nil {
			return nil, nil, err
		}

		for _, toArtifact := range pomDeps {
			var ok bool
			if _, ok = pomsByArtifact[fromArtifact]; !ok {
				// external dependency
				continue
			}
			var fromDeps []dependency
			if fromDeps, ok = depsByFrom[fromArtifact]; !ok {
				fromDeps = []dependency{}
				depsByFrom[fromArtifact] = fromDeps
			}
			var toDeps []dependency
			if toDeps, ok = depsByTo[toArtifact]; !ok {
				toDeps = []dependency{}
				depsByTo[toArtifact] = toDeps
			}
			dep := dependency{fromArtifact: fromArtifact, toArtifact: toArtifact}
			fromDeps = append(fromDeps, dep)
			toDeps = append(toDeps, dep)
		}
	}
	return depsByFrom, depsByTo, nil
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
