package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/cmd/columns"
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	"github.com/dmolesUC3/mrt-build-info/maven"
	. "github.com/dmolesUC3/mrt-build-info/shared"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	poms := &poms{}

	cmd := &cobra.Command{
		Use:   "poms",
		Short: "List Maven poms",
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
			return poms.List(server)
		},
	}
	cmd.Flags().BoolVarP(&poms.artifacts, "artifacts", "a", false, "list POM artifacts")
	cmd.Flags().BoolVarP(&poms.deps, "deps", "d", false, "list POM dependencies")
	cmd.Flags().BoolVarP(&maven.POMURLs, "pom-urls", "u", false, "list URL used to retrieve POM file")

	AddCommand(cmd)
}

type poms struct {
	artifacts bool
	deps bool
	errors    []error
}

//noinspection GoUnhandledErrorResult
func (p *poms) List(server jenkins.JenkinsServer) error {
	if Flags.Verbose {
		fmt.Fprintln(os.Stderr, "Retrieving jobs...")
	}
	allJobs, err := server.Jobs()
	if err != nil {
		return err
	}

	var jobs []jenkins.Job
	var poms []maven.Pom
	for _, j := range allJobs {
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
			jobs = append(jobs, j)
			poms = append(poms, p)
		}
	}

	cols := p.MakeTableColumns(jobs, poms)
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
		cols = append(cols, NewTableColumn("POM Blob URL", len(jobs), func(row int) string {
			url := poms[row].BlobURL()
			if url == nil {
				return columns.ValueUnknown
			}
			return url.String()
		}))
	}
	if p.artifacts {
		cols = append(cols, NewTableColumn("Artifacts", len(jobs), func(row int) string {
			return p.ArtifactInfo(jobs[row], poms[row])
		}))
	}
	if p.deps {
		cols = append(cols, NewTableColumn("Dependencies", len(jobs), func(row int) string {
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

func (p *poms) ArtifactInfo(job jenkins.Job, pom maven.Pom) string {
	artifact, err := pom.Artifact()
	if err != nil {
		p.errors = append(p.errors, err)
		return ""
	}
	artifactStr := artifact.String()
	if jenkins.IsParameterized(artifactStr) {
		var expanded []string
		artifactParams := jenkins.Parameters(artifactStr)
		mvnParamToVal, err := job.MavenParamToValue()
		if err != nil {
			p.errors = append(p.errors, err)
			return artifactStr
		}
		var missing []string
		var found []string
		for _, p := range artifactParams {
			current := len(expanded)
			if val, ok := mvnParamToVal[p]; ok {
				if strings.HasPrefix(val, "$") {
					for _, jp := range job.Parameters() {
						jpName := jp.Name()
						if val == "$"+jpName {
							found = append(found, p + " -> " + jpName)
							paramSub := "${" + p + "}"
							for _, v := range jp.Choices() {
								expanded = append(expanded, strings.ReplaceAll(artifactStr, paramSub, v))
							}
						}
					}
				}
			}
			if len(expanded) == current {
				missing = append(missing, p)
			}
		}
		if len(missing) > 0 {
			err = fmt.Errorf(
				"job %v: pom %v missing parameters: %v (found: %v)\n",
				job.Name(), pom.Path(), strings.Join(missing, ", "), strings.Join(found, ", "),
			)
			p.errors = append(p.errors, err)
		}
		return strings.Join(expanded, ", ")
	}
	return artifactStr
}
