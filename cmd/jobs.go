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
	cmd.Flags().BoolVarP(&jobs.parameters, "parameters", "p", false, "show build parameters")
	cmd.Flags().BoolVarP(&jobs.repo, "repositories", "r", false, "show repositories")

	cmd.Flags().BoolVar(&jobs.apiUrl, "api-url", false, "show Jenkins API URLs")
	cmd.Flags().BoolVar(&jobs.configXML, "config-xml", false, "show Jenkins config.xml URLs")
	cmd.Flags().BoolVar(&jobs.poms, "poms", false, "show POMs")

	AddCommand(cmd)
}

type jobs struct {
	apiUrl     bool
	artifacts  bool
	build      bool
	configXML  bool
	parameters bool
	poms       bool
	repo       bool
	errors     []error
}

//noinspection GoUnhandledErrorResult
func (j *jobs) List(server jenkins.JenkinsServer) error {
	jobs, err := server.Jobs()
	if err != nil {
		return err
	}

	columns := j.MakeTableColumns(jobs)
	if len(columns) == 1 {
		for row, rowCount := 0, columns[0].Rows(); row < rowCount; row++ {
			fmt.Println(columns[0].ValueAt(row))
		}
	} else {
		table := TableFrom(columns...)
		table.Print(os.Stdout, "\t")
	}

	if len(j.errors) > 0 {
		w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
		fmt.Fprintf(w, "%d errors:\n", len(j.errors))
		for i, err := range j.errors {
			fmt.Fprintf(w, "%d. %v\n", i+1, err)
		}
		w.Flush()
	}

	return nil
}

func (j *jobs) MakeTableColumns(jobs []jenkins.Job) []TableColumn {
	columns := []TableColumn{
		NewTableColumn("Job Name", len(jobs), func(row int) string {
			return jobs[row].Name()
		}),
	}
	if j.apiUrl {
		columns = append(columns, NewTableColumn(
			"API Url", len(jobs), func(row int) string {
				apiUrl := jobs[row].APIUrl()
				if apiUrl == nil {
					return valueUnknown
				}
				return apiUrl.String()
			}))
	}
	if j.configXML {
		columns = append(columns, NewTableColumn(
			"config.xml", len(jobs), func(row int) string {
				configUrl := jobs[row].ConfigUrl()
				if configUrl == nil {
					return valueUnknown
				}
				return configUrl.String()
			}))
	}
	if j.parameters {
		columns = append(columns, NewTableColumn(
			"Parameters", len(jobs), func(row int) string {
				var allParamInfo []string
				for _, p := range jobs[row].Parameters() {
					paramInfo := fmt.Sprintf("%v (%v)", p.Name(), strings.Join(p.Choices(), ", "))
					allParamInfo = append(allParamInfo, paramInfo)
				}
				return strings.Join(allParamInfo, ", ")
			}))
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
	if j.poms {
		columns = append(columns, NewTableColumn(
			"POMs", len(jobs), func(row int) string {
				poms, err := jobs[row].POMs()
				if err != nil {
					j.errors = append(j.errors, err)
					return valueUnknown
				}
				var allPomInfo []string
				for _, p := range poms {
					allPomInfo = append(allPomInfo, p.Path())
				}
				return strings.Join(allPomInfo, ", ")
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
	return columns
}
