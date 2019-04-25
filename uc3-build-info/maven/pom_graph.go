package maven

import (
	"fmt"
	. "github.com/CDLUC3/uc3-tools/mrt-build-info/shared"
	"os"
)

type PomGraph interface {
	Poms() []Pom
	DependenciesOf(pom Pom) (deps []Pom)
	DependenciesOn(pom Pom) (deps []Pom)
}

func NewPomGraph(agraph ArtifactGraph) PomGraph {
	return &pomGraph{artifactGraph: agraph}
}

type pomDep struct {
	fromPom Pom
	toPom   Pom
}

type pomGraph struct {
	artifactGraph ArtifactGraph

	allDeps    []pomDep
	depsByFrom map[Pom][]pomDep
	depsByTo   map[Pom][]pomDep
}

func (g *pomGraph) Poms() []Pom {
	return g.artifactGraph.Poms()
}

func (g *pomGraph) DependenciesOf(pom Pom) (deps []Pom) {
	_, depsByFrom, _ := g.deps()
	if depsFromPom, ok := depsByFrom[pom]; ok {
		for _, d := range depsFromPom {
			deps = append(deps, d.toPom)
		}
	}
	Deduplicate(PomsByLocation(deps), func(len int) { deps = deps[:len] })
	return deps
}

func (g *pomGraph) DependenciesOn(pom Pom) (deps []Pom) {
	_, _, depsByTo := g.deps()
	if depsToPom, ok := depsByTo[pom]; ok {
		for _, d := range depsToPom {
			deps = append(deps, d.fromPom)
		}
	}
	Deduplicate(PomsByLocation(deps), func(len int) { deps = deps[:len] })
	return deps
}

func (g *pomGraph) deps() (allDeps []pomDep, depsByFrom map[Pom][]pomDep, depsByTo map[Pom][]pomDep) {
	if g.allDeps == nil {
		allDeps = []pomDep{}
		depsByFrom = map[Pom][]pomDep{}
		depsByTo = map[Pom][]pomDep{}

		agraph := g.artifactGraph
		artifacts := agraph.Artifacts()
		//noinspection GoUnhandledErrorResult
		if Flags.Verbose {
			fmt.Fprintf(os.Stderr, "Determining pom dependencies for %d artifacts...", len(artifacts))
			defer func() { fmt.Fprintln(os.Stderr, "Done.") }()
		}
		for _, a := range artifacts {
			//noinspection GoUnhandledErrorResult
			if Flags.Verbose {
				fmt.Fprint(os.Stderr, ".")
			}
			fromArtifact := a
			fromPom := agraph.PomForArtifact(fromArtifact)
			for _, toArtifact := range agraph.DependenciesOf(fromArtifact) {
				toPom := agraph.PomForArtifact(toArtifact)

				var ok bool
				var fromDeps []pomDep
				if fromDeps, ok = depsByFrom[fromPom]; !ok {
					fromDeps = []pomDep{}
				}
				var toDeps []pomDep
				if toDeps, ok = depsByTo[toPom]; !ok {
					toDeps = []pomDep{}
				}

				dep := pomDep{fromPom: fromPom, toPom: toPom}
				allDeps = append(allDeps, dep)
				depsByFrom[fromPom] = append(fromDeps, dep)
				depsByTo[toPom] = append(toDeps, dep)
			}
		}
		g.allDeps = allDeps
		g.depsByFrom = depsByFrom
		g.depsByTo = depsByTo
	}
	return g.allDeps, g.depsByFrom, g.depsByTo
}
