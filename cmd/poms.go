package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	"github.com/dmolesUC3/mrt-build-info/maven"
	. "github.com/dmolesUC3/mrt-build-info/shared"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
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
	cmd.Flags().BoolVarP(&maven.POMURLs, "pom-urls", "u", false, "list URL used to retrieve POM file")

	AddCommand(cmd)
}

type poms struct {
	// // Deprecated
	// listUrls  bool
	artifacts bool
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
		jobPoms, err := j.POMs()
		if err != nil {
			p.errors = append(p.errors, err)
			continue
		}
		for _, p := range jobPoms {
			jobs = append(jobs, j)
			poms = append(poms, p)
		}
	}

	columns := p.MakeTableColumns(jobs, poms)
	table := TableFrom(columns...)
	table.Print(os.Stdout, "\t")

	if len(p.errors) > 0 {
		w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
		fmt.Fprintf(w, "%d errors:\n", len(p.errors))
		for i, err := range p.errors {
			fmt.Fprintf(w, "%d. %v\n", i+1, err)
		}
		w.Flush()
	}

	return nil
}

func (p *poms) MakeTableColumns(jobs []jenkins.Job, poms []maven.Pom) []TableColumn {
	if len(jobs) != len(poms) {
		panic(fmt.Errorf("mismatched jobs (%d) and poms(%d)", len(jobs), len(poms)))
	}
	rows := len(jobs)

	columns := []TableColumn{
		NewTableColumn("Job", rows, func(row int) string {
			return jobs[row].Name()
		}),
		NewTableColumn("POM", rows, func(row int) string {
			return poms[row].Path()
		}),
	}
	if maven.POMURLs {
		columns = append(columns, NewTableColumn("POM Blob URL", rows, func(row int) string {
			url := poms[row].BlobURL()
			if url == nil {
				return valueUnknown
			}
			return url.String()
		}))
	}
	if p.artifacts {
		columns = append(columns, NewTableColumn("Artifacts", rows, func(row int) string {
			return p.ArtifactInfo(jobs[row], poms[row])
		}))
	}

	return columns
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
