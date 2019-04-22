package cmd

import (
	"fmt"
	. "github.com/dmolesUC3/mrt-build-info/dependencies"
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	. "github.com/dmolesUC3/mrt-build-info/shared"
	"github.com/spf13/cobra"
	"os"
)

// TODO: rewrite to use new graph types

func init() {
	deps := &deps{}

	cmd := &cobra.Command{
		Use:   "deps",
		Short: "List internal Maven dependencies",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			server, err := ServerFrom(args)
			if err != nil {
				return err
			}
			return deps.List(server)
		},
	}

	cmd.Flags().BoolVarP(&deps.jobs, "jobs", "j", false, "list Jenkins job dependencies")
	cmd.Flags().BoolVarP(&deps.poms, "poms", "p", false, "list Maven pom dependencies")
	cmd.Flags().BoolVarP(&deps.artifacts, "artifacts", "a", false, "list Maven artifact dependencies")
	cmd.Flags().BoolVar(&deps.all, "all", false, "list all depdendencies")

	AddCommand(cmd)
}

type deps struct {
	jobs      bool
	poms      bool
	artifacts bool
	all       bool
	errors    []error
}

func (d *deps) unselected() bool {
	return !(d.jobs || d.poms || d.artifacts)
}

//noinspection GoUnhandledErrorResult
func (d *deps) List(server jenkins.JenkinsServer) error {
	if d.all { // TODO: see if cobra provides a more elegant --all
		d.jobs = true
		d.poms = true
		d.artifacts = true
	}
	if d.unselected() {
		return fmt.Errorf("at least one of --jobs, --poms, or --artifacts (or --all) is required")
	}
	jobs, err := server.Jobs()
	if err != nil {
		return err
	}

	jgraph := jenkins.NewJobGraph(jobs)

	var titles []string
	var tables []func() Table

	if d.jobs {
		titles = append(titles, "Jenkins job dependencies")
		tables = append(tables, func() Table { return JobsTable(jgraph) })
	}
	if d.poms {
		titles = append(titles, "Maven pom dependencies")
		tables = append(tables, func() Table {
			table, errs := PomsTable(jgraph)
			if len(errs) > 0 {
				d.errors = append(d.errors, errs...)
			}
			return table
		})
	}
	if d.artifacts {
		titles = append(titles, "Maven artifact dependencies")
		tables = append(tables, func() Table {
			table, errs := ArtifactsTable(jgraph)
			if len(errs) > 0 {
				d.errors = append(d.errors, errs...)
			}
			return table
		})
	}
	for i, title := range titles {
		//noinspection GoNilness
		table := tables[i]

		fmt.Println(title)
		fmt.Println()

		table().Print(os.Stdout, "\t")
		if i+1 < len(titles) {
			fmt.Println()
		}
	}
	return nil
}
