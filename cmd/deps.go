package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/cmd/columns"
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	"github.com/dmolesUC3/mrt-build-info/maven"
	. "github.com/dmolesUC3/mrt-build-info/shared"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	deps := &deps{}

	cmd := &cobra.Command{
		Use:   "deps",
		Short: "List internal Maven dependencies",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var server jenkins.JenkinsServer
			if len(args) == 0 {
				server = jenkins.DefaultServer()
			} else {
				server, err = jenkins.ServerFromUrl(args[0])
				if err != nil {
					return err
				}
			}
			return deps.List(server)
		},
	}

	AddCommand(cmd)
}

type deps struct {
	errors []error
}

// TODO: clean this up, move most of it somewhere (Jenkins package?)
//noinspection GoUnhandledErrorResult
func (d *deps) List(server jenkins.JenkinsServer) error {
	allJobs, err := server.Jobs()
	if err != nil {
		return err
	}

	if Flags.Verbose {
		fmt.Fprintln(os.Stderr, "Retreiving POMs...")
	}
	var jobs []jenkins.Job
	var poms []maven.Pom
	jobsByPom := map[maven.Pom]jenkins.Job{}
	for _, j := range allJobs {
		if Flags.Job != "" && j.Name() != Flags.Job {
			continue
		}
		jobPoms, errs := j.POMs()
		if len(errs) > 0 {
			d.errors = append(d.errors, errs...)
		}
		if len(jobPoms) == 0 {
			d.errors = append(d.errors, fmt.Errorf("no POMs found for job %v", j.Name()))
		}
		for _, p := range jobPoms {
			jobs = append(jobs, j)
			poms = append(poms, p)
			jobsByPom[p] = j
		}
	}

	graph, errors := maven.NewGraph(poms)
	d.errors = append(d.errors, errors...)

	artifacts := graph.SortedArtifacts()

	allInfo := map[maven.Artifact]*artifactInfo{}
	infoFor := func(artifact maven.Artifact) *artifactInfo {
		var info *artifactInfo
		var ok bool
		if info, ok = allInfo[artifact]; !ok {
			pom := graph.PomForArtifact(artifact)
			info = &artifactInfo{job: jobsByPom[pom], pom: pom, artifact: artifact}
			allInfo[artifact] = info
		}
		return info
	}

	var rows []depRow

	// TODO: table
	// TODO: job deps vs pom deps
	for _, art := range artifacts {
		artInfo := infoFor(art)

		requires := graph.DependenciesOf(art)
		requiredBy := graph.DependenciesOn(art)

		artRows := 1
		if len(requires) > artRows {
			artRows = len(requires)
		}
		if len(requiredBy) > artRows {
			artRows = len(requiredBy)
		}

		for i := 0; i < artRows; i++ {
			row := depRow{artifact: artInfo}
			if i < len(requires) {
				row.requires = infoFor(requires[i])
			}
			if i < len(requiredBy) {
				row.requiredBy = infoFor(requiredBy[i])
			}
			rows = append(rows, row)
		}
	}

	cols := makeTableColumns(rows)
	table := TableFrom(cols...)
	table.Print(os.Stdout, "\t")

	PrintErrors(d.errors)
	return nil
}

type artifactInfo struct {
	job      jenkins.Job
	pom      maven.Pom
	artifact maven.Artifact
}

func (a *artifactInfo) Job() jenkins.Job {
	if a == nil {
		return nil
	}
	return a.job
}

func (a *artifactInfo) Pom() maven.Pom {
	if a == nil {
		return nil
	}
	return a.pom
}

func (a *artifactInfo) Artifact() maven.Artifact {
	if a == nil {
		return nil
	}
	return a.artifact
}

type depRow struct {
	artifact   *artifactInfo
	requires   *artifactInfo
	requiredBy *artifactInfo
}

func makeTableColumns(rows []depRow) []TableColumn {
	return []TableColumn{
		columns.Jobs(func(row int) jenkins.Job { return rows[row].artifact.Job() }, len(rows)),
		columns.Poms(func(row int) maven.Pom { return rows[row].artifact.Pom() }, len(rows)),
		columns.Artifacts(func(row int) maven.Artifact { return rows[row].artifact.Artifact() }, len(rows)),

		columns.Jobs(func(row int) jenkins.Job { return rows[row].requires.Job() }, len(rows)),
		columns.Poms(func(row int) maven.Pom { return rows[row].requires.Pom() }, len(rows)),
		columns.Artifacts(func(row int) maven.Artifact { return rows[row].requires.Artifact() }, len(rows)),

		columns.Jobs(func(row int) jenkins.Job { return rows[row].requiredBy.Job() }, len(rows)),
		columns.Poms(func(row int) maven.Pom { return rows[row].requiredBy.Pom() }, len(rows)),
		columns.Artifacts(func(row int) maven.Artifact { return rows[row].requiredBy.Artifact() }, len(rows)),
	}
}
