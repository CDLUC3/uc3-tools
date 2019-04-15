package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/git"
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	"github.com/dmolesUC3/mrt-build-info/maven"
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

	rootCmd.AddCommand(cmd)
}

type poms struct {
	token     string
	logErrors bool
	listUrls  bool
	fullSHA bool
}

func (p *poms) List(server jenkins.JenkinsServer) error {
	jobs, err := server.Jobs()
	if err != nil {
		return err
	}

	for _, job := range jobs {
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
		repository := git.NewRepository(owner, repo, sha1, p.token)
		pomEntries, err := repository.Find("pom.xml$", git.Blob)
		if err != nil {
			if p.logErrors {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
			continue
		}
		for _, e := range pomEntries {
			pom, err := maven.PomFromEntry(e)
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
			pomInfo := fmt.Sprintf("%v\t%v\t%v", artifact, pom.Repository(), pom.Path())
			if p.listUrls {
				pomInfo = fmt.Sprintf("%v\t%v", pomInfo, git.WebUrlForEntry(e, p.fullSHA))
			}
			fmt.Println(pomInfo)
		}
	}
	return nil
}
