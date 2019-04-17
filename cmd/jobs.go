package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	. "github.com/dmolesUC3/mrt-build-info/shared"
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
	cmd.Flags().BoolVarP(&jobs.artifacts, "artifacts", "a", false, "list artifacts")
	cmd.Flags().BoolVarP(&jobs.build, "build", "b", false, "show info for last successful build")
	cmd.Flags().BoolVarP(&jobs.repo, "repositories", "r", false, "list repositories")

	AddCommand(cmd)
}

type jobs struct {
	artifacts bool
	build     bool
	repo      bool
}

func (j *jobs) nameOnly() bool {
	return !(j.artifacts || j.build || j.repo)
}

func (j *jobs) List(server jenkins.JenkinsServer) error {
	jobs, err := server.Jobs()
	if err != nil {
		return err
	}

	// TODO: some kind of column model that makes this less hacky
	if !j.nameOnly() {
		var fields []string
		fields = append(fields, "Job Name")
		if j.repo {
			fields = append(fields, "Repository")
		}
		if j.build {
			fields = append(fields, "Build")
			if Flags.Verbose {
				fields = append(fields, "SHA Hash")
			}
		}
		if j.artifacts {
			fields = append(fields, "Artifacts")
		}
		fmt.Println(strings.Join(fields, "\t"))
	}

	for _, job := range jobs {
		if Flags.Job != "" && Flags.Job != job.Name() {
			continue
		}
		if err := j.printJob(job); err != nil {
			if Flags.LogErrors {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}
	return nil
}

func (j *jobs) printJob(job jenkins.Job) error {
	name := job.Name()
	if j.nameOnly() {
		fmt.Println(name)
		return nil
	}
	b, err := job.LastSuccess()
	if err != nil {
		return err
	}

	var fields []string
	fields = append(fields, name)

	// TODO: some kind of column model that makes this less hacky
	if j.repo {
		scmUrl, err := b.SCMUrl()
		if err != nil {
			if Flags.LogErrors {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
			fields = append(fields, "")
		}
		fields = append(fields, scmUrl)
	}

	if j.build {
		fields = append(fields, fmt.Sprintf("%d", b.BuildNumber()))
		if Flags.Verbose {
			sha1, err := b.SHA1()
			if err != nil {
				if Flags.LogErrors {
					_, _ = fmt.Fprintln(os.Stderr, err.Error())
				}
				fields = append(fields, "")
			} else {
				fields = append(fields, sha1.String())
			}
		}
	}

	if j.artifacts {
		artifacts, err := b.Artifacts()
		if err != nil {
			if Flags.LogErrors {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
			fields = append(fields, "")
		} else {
			var allArtifactInfo []string
			for _, a := range artifacts {
				allArtifactInfo = append(allArtifactInfo, a.String())
			}
			fields = append(fields, strings.Join(allArtifactInfo, ", "))
		}
	}

	fmt.Println(strings.Join(fields, "\t"))

	return nil
}
