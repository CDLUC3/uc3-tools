package maven

import (
	"fmt"
	. "github.com/CDLUC3/uc3-tools/uc3-build-info/shared"
	"os"
	"sort"
)

// ------------------------------------------------------------
// ArtifactGraph

type ArtifactGraph interface {
	Artifacts() []Artifact
	SortedArtifacts() []Artifact
	DependenciesOf(artifact Artifact) (deps []Artifact)
	DependenciesOn(artifact Artifact) (deps []Artifact)
	PomForArtifact(artifact Artifact) Pom
	Poms() []Pom
}

func NewArtifactGraph(poms []Pom) (ArtifactGraph, []error) {
	//noinspection GoUnhandledErrorResult
	if Flags.Verbose {
		fmt.Fprintf(os.Stderr, "Determining artifact dependencies for %d poms...\n", len(poms))
	}

	var errors []error
	artifacts, pomsByArtifact, errs1 := artifactsFromPoms(poms)
	errors = append(errors, errs1...)

	depsByFrom, depsByTo, errs2 := dependencies(poms, pomsByArtifact)
	errors = append(errors, errs2...)

	g := artifactGraph{poms: poms, artifacts: artifacts, pomsByArtifact: pomsByArtifact, depsByFrom: depsByFrom, depsByTo: depsByTo}
	return &g, errors
}

// ------------------------------------------------------------
// Unexported symbols

type artifactDep struct {
	fromArtifact Artifact
	toArtifact   Artifact
}

type artifactGraph struct {
	poms []Pom
	artifacts      []Artifact
	pomsByArtifact map[Artifact]Pom
	depsByFrom     map[Artifact][]artifactDep
	depsByTo       map[Artifact][]artifactDep

	sortedArtifacts []Artifact
}

func (g *artifactGraph) Poms() []Pom {
	return g.poms
}

func (g *artifactGraph) DependenciesOf(artifact Artifact) (deps []Artifact) {
	if depsFromArtifact, ok := g.depsByFrom[artifact]; ok {
		for _, d := range depsFromArtifact {
			deps = append(deps, d.toArtifact)
		}
	}
	Deduplicate(ArtifactsByString(deps), func(len int) { deps = deps[:len] })
	return deps
}

func (g *artifactGraph) DependenciesOn(artifact Artifact) (deps []Artifact) {
	if depsToArtifact, ok := g.depsByTo[artifact]; ok {
		for _, d := range depsToArtifact {
			deps = append(deps, d.fromArtifact)
		}
	}
	Deduplicate(ArtifactsByString(deps), func(len int) { deps = deps[:len] })
	return deps
}

func (g *artifactGraph) Artifacts() []Artifact {
	return g.artifacts
}

func (g *artifactGraph) SortedArtifacts() []Artifact {
	if g.sortedArtifacts == nil {
		//noinspection GoUnhandledErrorResult
		if Flags.Verbose {
			fmt.Fprintf(os.Stderr, "Sorting %d artifacts...", len(g.artifacts))
			defer func() { fmt.Fprintln(os.Stderr, "Done.") }()
		}
		g.sortedArtifacts = newTopoSort(g).sortedArtifacts()
	}
	return g.sortedArtifacts
}

func (g *artifactGraph) PomForArtifact(artifact Artifact) Pom {
	return g.pomsByArtifact[artifact]
}

func artifactsFromPoms(poms []Pom) ([]Artifact, map[Artifact]Pom, []error) {
	//noinspection GoUnhandledErrorResult
	if Flags.Verbose {
		fmt.Fprintf(os.Stderr, "Identifying artifacts for %d poms...", len(poms))
		defer func() { fmt.Fprintln(os.Stderr, "Done.") }()
	}

	var errors []error

	var allArtifacts []Artifact
	a2p := map[Artifact]Pom{}
	for _, pom := range poms {
		//noinspection GoUnhandledErrorResult
		if Flags.Verbose {
			fmt.Fprint(os.Stderr, ".")
		}
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

func dependencies(poms []Pom, pomsByArtifact map[Artifact]Pom) (depsByFrom map[Artifact][]artifactDep, depsByTo map[Artifact][]artifactDep, err []error) {
	//noinspection GoUnhandledErrorResult
	if Flags.Verbose {
		fmt.Fprintf(os.Stderr, "Determining artifact dependencies for %d poms...", len(poms))
		defer func() { fmt.Fprintln(os.Stderr, "Done.") }()
	}

	var errors []error

	depsByFrom = map[Artifact][]artifactDep{}
	depsByTo = map[Artifact][]artifactDep{}
	for _, pom := range poms {
		//noinspection GoUnhandledErrorResult
		if Flags.Verbose {
			fmt.Fprint(os.Stderr, ".")
		}
		fromArtifact, err := pom.Artifact()
		if err != nil {
			errors = append(errors, err)
			continue
		}
		pomDeps, depErrs := pom.Dependencies()
		errors = append(errors, depErrs...)

		for _, toArtifact := range pomDeps {
			if fromArtifact == toArtifact {
				continue
			}
			var ok bool
			if _, ok = pomsByArtifact[toArtifact]; !ok {
				// third-party artifactDep
				continue
			}
			var fromDeps []artifactDep
			if fromDeps, ok = depsByFrom[fromArtifact]; !ok {
				fromDeps = []artifactDep{}
			}
			var toDeps []artifactDep
			if toDeps, ok = depsByTo[toArtifact]; !ok {
				toDeps = []artifactDep{}
			}
			dep := artifactDep{fromArtifact: fromArtifact, toArtifact: toArtifact}
			depsByFrom[fromArtifact] = append(fromDeps, dep)
			depsByTo[toArtifact] = append(toDeps, dep)
		}
	}
	return depsByFrom, depsByTo, errors
}

// ------------------------------------------------------------
// Topological sorter

type topoSort struct {
	g        *artifactGraph
	sorted   []Artifact
	visiting map[Artifact]bool
	visited  map[Artifact]bool
}

func newTopoSort(g *artifactGraph) *topoSort {
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
	//noinspection GoUnhandledErrorResult
	if Flags.Verbose {
		fmt.Fprint(os.Stderr, ".")
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
