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
	poms := &poms{}

	cmd := &cobra.Command{
		Use:   "poms",
		Short: "List Maven poms",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			server, err := ServerFrom(args)
			if err != nil {
				return err
			}
			return poms.List(server)
		},
	}
	cmd.Flags().BoolVarP(&poms.artifacts, "artifacts", "a", false, "list POM artifacts")
	cmd.Flags().BoolVarP(&poms.deps, "deps", "d", false, "list POM dependencies")
	cmd.Flags().BoolVarP(&maven.POMURLs, "pom-urls", "u", false, "list URL used to retrieve POM file")
	cmd.Flags().StringVarP(&Flags.Job, "job", "j", "", "show info only for specified job")

	AddCommand(cmd)
}

type poms struct {
	artifacts bool
	deps bool
	errors    []error
}

//noinspection GoUnhandledErrorResult
func (p *poms) List(server jenkins.JenkinsServer) error {
	jobs, err := server.Jobs()
	if err != nil {
		return err
	}
	jgraph := jenkins.NewJobGraph(jobs)
	jgraph.PomGraph()


	var pomJobs []jenkins.Job
	var poms []maven.Pom
	for _, j := range jobs {
		if Flags.Job != "" && j.Name() != Flags.Job {
			continue
		}
		jobPoms, errs := j.POMs()
		if len(errs) > 0 {
			p.errors = append(p.errors, errs...)
		}
		if len(jobPoms) == 0 {
			p.errors = append(p.errors, fmt.Errorf("no POMs found for job %v", j.Name()))
		}
		for _, p := range jobPoms {
			pomJobs = append(pomJobs, j)
			poms = append(poms, p)
		}
	}

	cols := p.MakeTableColumns(pomJobs, poms)
	table := TableFrom(cols...)
	table.Print(os.Stdout, "\t")

	PrintErrors(p.errors)
	return nil
}


func (p *poms) MakeTableColumns(jobs []jenkins.Job, poms []maven.Pom) []TableColumn {
	if len(jobs) != len(poms) {
		panic(fmt.Errorf("mismatched jobs (%d) and poms(%d)", len(jobs), len(poms)))
	}
	cols := []TableColumn{columns.Job(jobs), columns.Pom(poms), }
	if maven.POMURLs {
		cols = append(cols, NewTableColumn("POM Blob URL", len(poms), func(row int) string {
			url := poms[row].BlobURL()
			if url == nil {
				return columns.ValueUnknown
			}
			return url.String()
		}))
	}
	if p.artifacts {
		cols = append(cols, NewTableColumn("Artifacts", len(poms), func(row int) string {
			pom := poms[row]
			artifact, err := pom.Artifact()
			if err != nil {
				p.errors = append(p.errors, err)
				return ""
			}
			return artifact.String()
		}))
	}
	if p.deps {
		cols = append(cols, NewTableColumn("Dependencies", len(poms), func(row int) string {
			deps, errs := poms[row].Dependencies()
			if errs != nil {
				p.errors = append(p.errors, errs...)
				return ""
			}
			return maven.ArtifactsByString(deps).String()
		}))
	}

	return cols
}
