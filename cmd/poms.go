package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/git"
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	"github.com/dmolesUC3/mrt-build-info/maven"
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
	cmd.Flags().StringVarP(&poms.token, "token", "t", "", "GitHub API token (https://github.com/settings/tokens)")
	cmd.Flags().BoolVarP(&poms.logErrors, "log-errors", "l", false, "log non-fatal errors to stderr")
	cmd.Flags().BoolVarP(&poms.listUrls, "list-urls", "u", false, "list URL used to retrieve POM file")
	cmd.Flags().BoolVarP(&poms.fullSHA, "full-sha", "f", false, "don't abbreviate SHA hashes in URLs")
	cmd.Flags().BoolVarP(&poms.verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(cmd)
}

type poms struct {
	token     string
	logErrors bool
	listUrls  bool
	fullSHA   bool
	verbose   bool
}

func (p *poms) List(server jenkins.JenkinsServer) error {
	jobs, err := server.Jobs()
	if err != nil {
		return err
	}

	// TODO: extract some more methods, clean this up (continue -> return err?)
	for i, job := range jobs {
		if p.verbose {
			fmt.Printf("Job: %v (%d of %d)\n", job.Name(), i+1, len(jobs))
		}
		build, err := job.LastSuccess()
		if err != nil {
			if p.logErrors {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
			continue
		}
		owner, repo, sha1, err := build.Commit()
		if err != nil {
			if p.logErrors {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
			continue
		}
		repository, err := git.GetRepository(owner, repo, sha1, p.token)
		if err != nil {
			if p.logErrors {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
			continue
		}

		pomEntries, err := repository.Find("pom.xml$", git.Blob)
		if err != nil {
			if p.logErrors {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
			continue
		}
		for _, entry := range pomEntries {
			pom, err := maven.PomFromEntry(entry)
			if err != nil {
				if p.logErrors {
					_, _ = fmt.Fprintln(os.Stderr, err.Error())
				}
				continue
			}
			artifact, err := pom.Artifact()
			if err != nil {
				if p.logErrors {
					_, _ = fmt.Fprintln(os.Stderr, err.Error())
				}
				continue
			}
			artifactStr := artifact.String()
			if strings.Contains(artifactStr, "$") {
				err = p.printParameterizedPomInfo(artifactStr, pom, entry, job.Name(), job.Parameters())
				if p.logErrors {
					_, _ = fmt.Fprintln(os.Stderr, err.Error())
				}
			} else {
				p.printPomInfo(artifactStr, pom, entry)
			}
		}
	}
	return nil
}

func (p *poms) printParameterizedPomInfo(artifactStr string, pom maven.Pom, entry git.Entry, jobName string, params []jenkins.Parameter) error {
	didPrint := false
	for _, p := range params {
		paramSub := "${" + p.Name() + "}"
		if strings.Contains(artifactStr, paramSub) {
			for _, v := range p.Choices() {
				artifactStr2 := strings.ReplaceAll(artifactStr, paramSub, v)
				print(artifactStr2, pom, entry)
				didPrint = true
			}
		}
	}
	if didPrint {
		return nil
	}
	return fmt.Errorf("no matching parameter in job %#v for substitution in %#v", jobName, artifactStr)
}

func (p *poms) printPomInfo(artifactStr string, pom maven.Pom, entry git.Entry) {
	pomInfo := fmt.Sprintf("%v\t%v\t%v", artifactStr, pom.Repository(), pom.Path())
	if p.listUrls {
		pomInfo = fmt.Sprintf("%v\t%v", pomInfo, git.WebUrlForEntry(entry, p.fullSHA))
	}
	fmt.Println(pomInfo)
}
