package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	"github.com/dmolesUC3/mrt-build-info/maven"
	. "github.com/dmolesUC3/mrt-build-info/shared"
	"github.com/spf13/cobra"
	"os"
	"strings"
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

//noinspection GoUnhandledErrorResult
func (d *deps) List(server jenkins.JenkinsServer) error {
	if Flags.Verbose {
		fmt.Fprintln(os.Stderr, "Retrieving jobs...")
	}
	allJobs, err := server.Jobs()
	if err != nil {
		return err
	}

	if Flags.Verbose {
		fmt.Fprintln(os.Stderr, "Retreiving POMs...")
	}
	jobsByPom := map[maven.Pom]jenkins.Job{}
	var poms []maven.Pom
	for _, j := range allJobs {
		if Flags.Job != "" && j.Name() != Flags.Job {
			continue
		}
		jobPoms, err := j.POMs()
		if err != nil {
			d.errors = append(d.errors, err)
			continue
		}
		for _, p := range jobPoms {
			jobsByPom[p] = j
			poms = append(poms, p)
		}
	}

	if Flags.Verbose {
		fmt.Fprintln(os.Stderr, "Building graph...")
	}
	graph, errors := maven.NewGraph(poms)
	d.errors = append(d.errors, errors...)

	if Flags.Verbose {
		fmt.Fprintln(os.Stderr, "Sorting POMs...")
	}
	artifacts := graph.SortedArtifacts()

	// TODO: prettier output
	// TODO: job deps vs pom deps
	for _, artifact := range artifacts {
		fmt.Println(artifact.String())

		requires := graph.DependenciesOf(artifact)
		requiresStr := make([]string, len(requires))
		for i, r := range requires {
			requiresStr[i] = r.String()
		}

		// TODO: figure out why this is always empty
		requiredBy := graph.DependenciesOn(artifact)
		requiredByStr := make([]string, len(requiredBy))
		for i, r := range requiredBy {
			requiredByStr[i] = r.String()
		}

		fmt.Printf("\tRequires:    %v\n", strings.Join(requiresStr, ", "))
		fmt.Printf("\tRequired by: %v\n", strings.Join(requiredByStr, ", "))
	}

	PrintErrors(d.errors)
	return nil
}
