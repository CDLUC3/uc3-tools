package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	. "github.com/dmolesUC3/mrt-build-info/shared"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

const valueUnknown = "(unknown)"

func init() {
	jobs := &jobs{}
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "List Jenkins jobs",
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
			return jobs.List(server)
		},
	}
	cmd.Flags().BoolVarP(&jobs.artifacts, "artifacts", "a", false, "show artifacts from last successful build")
	cmd.Flags().BoolVarP(&jobs.build, "build", "b", false, "show info for last successful build")
	cmd.Flags().BoolVarP(&jobs.repo, "repositories", "r", false, "show repositories")

	AddCommand(cmd)
}

type jobs struct {
	artifacts bool
	build     bool
	repo      bool
	errors    []error
}

func (j *jobs) nameOnly() bool {
	return !(j.artifacts || j.build || j.repo)
}

func (j *jobs) List(server jenkins.JenkinsServer) error {
	jobs, err := server.Jobs()
	if err != nil {
		return err
	}

	if j.nameOnly() {
		for row := 0; row < len(jobs); row++ {
			fmt.Println(jobs[row].Name())
		}
	} else {
		j.printTable(jobs)
	}

	return nil
}

//noinspection GoUnhandledErrorResult
func (j *jobs) printTable(jobs []jenkins.Job) {
	columns := []TableColumn{
		NewTableColumn("Job Name", len(jobs), func(row int) string {
			return jobs[row].Name()
		}),
	}
	if j.repo {
		columns = append(columns, NewTableColumn(
			"Repository", len(jobs), func(row int) string {
				scmUrl, err := jobs[row].SCMUrl()
				if err != nil {
					j.errors = append(j.errors, err)
					return valueUnknown
				}
				return scmUrl
			}))
	}
	if j.build {
		columns = append(columns, NewTableColumn(
			"Last Success", len(jobs), func(row int) string {
				b, err := jobs[row].LastSuccess()
				if err != nil {
					j.errors = append(j.errors, err)
					return valueUnknown
				}
				return fmt.Sprintf("%d", b.BuildNumber())
			}))
		columns = append(columns, NewTableColumn(
			"SHA Hash", len(jobs), func(row int) string {
				b, err := jobs[row].LastSuccess()
				if err != nil {
					j.errors = append(j.errors, err)
					return valueUnknown
				}
				sha1, err := b.SHA1()
				if err != nil {
					j.errors = append(j.errors, err)
					return valueUnknown
				}
				return sha1.String()
			}))
	}
	if j.artifacts {
		columns = append(columns, NewTableColumn(
			"Last Artifacts", len(jobs), func(row int) string {
				b, err := jobs[row].LastSuccess()
				if err != nil {
					j.errors = append(j.errors, err)
					return valueUnknown
				}
				artifacts, err := b.Artifacts()
				if len(artifacts) == 0 {
					return "(no artifacts)"
				}
				var allArtifactInfo []string
				for _, a := range artifacts {
					allArtifactInfo = append(allArtifactInfo, a.String())
				}
				return strings.Join(allArtifactInfo, ", ")
			}))
	}
	table := TableFrom(columns...)
	table.Print(os.Stdout, "\t")

	if len(j.errors) > 0 {
		w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
		fmt.Fprintf(w, "%d errors:\n", len(j.errors))
		for i, err := range j.errors {
			fmt.Fprintf(w, "%d. %v\n", i+1, err)
		}
		w.Flush()
	}
}
