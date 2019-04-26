package cmd

import (
	"fmt"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/cmd/columns"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/jenkins"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/maven"
	. "github.com/CDLUC3/uc3-tools/uc3-build-info/shared"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	jobs := &jobs{}
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "List Jenkins jobs",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			server, err := ServerFrom(args)
			if err != nil {
				return err
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
	cmd.Flags().BoolVar(&jobs.goals, "goals", false, "show Maven goals")
	cmd.Flags().StringVarP(&Flags.Job, "job", "j", "", "show info only for specified job")

	AddCommand(cmd)
}

type jobs struct {
	apiUrl     bool
	artifacts  bool
	build      bool
	configXML  bool
	goals      bool
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

	cols := j.MakeTableColumns(jobs)
	if len(cols) == 1 {
		for row, rowCount := 0, cols[0].Rows(); row < rowCount; row++ {
			fmt.Println(cols[0].ValueAt(row))
		}
	} else {
		table := TableFrom(cols...)
		table.Print(os.Stdout, "\t")
	}

	PrintErrors(j.errors)
	return nil
}

func (j *jobs) MakeTableColumns(jobs []jenkins.Job) []TableColumn {
	rows := len(jobs)
	cols := []TableColumn{columns.Job(jobs)}
	if j.apiUrl {
		cols = append(cols, NewTableColumn(
			"API Url", rows, func(row int) string {
				apiUrl := jobs[row].APIUrl()
				if apiUrl == nil {
					return columns.ValueUnknown
				}
				return apiUrl.String()
			}))
	}
	if j.configXML {
		cols = append(cols, NewTableColumn(
			"config.xml", rows, func(row int) string {
				configUrl := jobs[row].ConfigUrl()
				if configUrl == nil {
					return columns.ValueUnknown
				}
				return configUrl.String()
			}))
	}
	if j.parameters {
		cols = append(cols, NewTableColumn(
			"Parameters", rows, func(row int) string {
				var allParamInfo []string
				for _, p := range jobs[row].Parameters() {
					paramInfo := fmt.Sprintf("%v (%v)", p.Name(), strings.Join(p.Choices(), ", "))
					allParamInfo = append(allParamInfo, paramInfo)
				}
				return strings.Join(allParamInfo, ", ")
			}))
	}
	if j.repo {
		cols = append(cols, NewTableColumn(
			"Repository", rows, func(row int) string {
				scmUrl, err := jobs[row].SCMUrl()
				if err != nil {
					j.errors = append(j.errors, err)
					return columns.ValueUnknown
				}
				return scmUrl
			}))
	}
	if j.poms {
		cols = append(cols, NewTableColumn(
			"POMs", rows, func(row int) string {
				poms, errs := jobs[row].POMs()
				if len(errs) > 0 {
					j.errors = append(j.errors, errs...)
				}
				if len(poms) == 0 {
					return columns.ValueUnknown
				}
				var allPomInfo []string
				for _, p := range poms {
					allPomInfo = append(allPomInfo, p.Path())
				}
				return strings.Join(allPomInfo, ", ")
			}))
	}
	if j.goals {
		cols = append(cols, NewTableColumn(
			"Goals", rows, func(row int) string {
				config, err := jobs[row].Config()
				if err != nil {
					j.errors = append(j.errors, err)
					return columns.ValueUnknown
				}
				return config.Goals()
			}))
	}
	if j.build {
		cols = append(cols, NewTableColumn(
			"Last Success", rows, func(row int) string {
				b, err := jobs[row].LastSuccess()
				if err != nil {
					j.errors = append(j.errors, err)
					return columns.ValueUnknown
				}
				return fmt.Sprintf("%d", b.BuildNumber())
			}))
		cols = append(cols, NewTableColumn(
			"SHA Hash", rows, func(row int) string {
				b, err := jobs[row].LastSuccess()
				if err != nil {
					j.errors = append(j.errors, err)
					return columns.ValueUnknown
				}
				sha1, err := b.SHA1()
				if err != nil {
					j.errors = append(j.errors, err)
					return columns.ValueUnknown
				}
				return sha1.String()
			}))
	}
	if j.artifacts {
		cols = append(cols, NewTableColumn(
			"Last Artifacts", rows, func(row int) string {
				b, err := jobs[row].LastSuccess()
				if err != nil {
					j.errors = append(j.errors, err)
					return columns.ValueUnknown
				}
				artifacts, err := b.Artifacts()
				if len(artifacts) == 0 {
					return "(no artifacts)"
				}
				return maven.ArtifactsByString(artifacts).String()
			}))
	}
	return cols
}
