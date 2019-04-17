package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/git"
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
	cmd.Flags().BoolVarP(&maven.POMURLs, "pom-urls", "u", false, "list URL used to retrieve POM file")

	AddCommand(cmd)
}

type poms struct {
	listUrls  bool
	token     string
}

// TODO:
//   - Don't just iterate over every pom in every repo; read the builds
//   - Move print logic into domain objects
//   - Get groupId etc. from parent POMs

func (p *poms) List(server jenkins.JenkinsServer) error {
	jobs, err := server.Jobs()
	if err != nil {
		return err
	}

	job := Flags.Job
	if job == "" {
		p.printAllJobs(jobs)
		return nil
	}

	for _, j := range jobs {
		if j.Name() == job {
			return p.printOneJob(j)
		}
	}
	return fmt.Errorf("no such job: %#v", job)
}

func (p *poms) printAllJobs(jobs []jenkins.Job) {
	for i, j := range jobs {
		if Flags.Verbose {
			fmt.Printf("%v (job %d of %d):\n", j.Name(), i+1, len(jobs))
		}
		err := p.printOneJob(j)
		if err != nil  {
			if Flags.Verbose {
				_, _ = fmt.Fprint(os.Stderr, "\t")
			}
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

func (p *poms) printOneJob(j jenkins.Job) error {
	j1 := job{*p, j}
	pomEntries, err := j1.pomEntries()
	if err != nil {
		return err
	}
	for i, pomEntry := range pomEntries {
		if Flags.Verbose {
			fmt.Printf("\t%v: %v (pom %d of %d)\n\t\t", j.Name(), pomEntry.Path(), i+1, len(pomEntries))
		}
		err = j1.printPomEntry(pomEntry)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	return nil
}

type job struct {
	poms // TODO: don't do this, just explicitly include fields
	jenkins.Job
}

func (j *job) repo() (git.Repository, error) {
	build, err := j.LastSuccess()
	if err != nil {
		return nil, fmt.Errorf("can't determine repository for job %v: %v", j.Name(), err)
	}
	owner, repoName, sha1, err := build.Commit()
	if err != nil {
		return nil, fmt.Errorf("can't determine repository for job %v: %v", j.Name(), err)
	}
	repo, err := git.GetRepository(owner, repoName, sha1)
	if err != nil {
		return nil, fmt.Errorf("can't determine repository for job %v: %v", j.Name(), err)
	}
	return repo, nil
}

func (j *job) pomEntries() ([]git.Entry, error) {
	repository, err := j.repo()
	if err != nil {
		return nil, err
	}
	return repository.Find("pom.xml$", git.Blob)
}

func (j *job) printPomEntry(entry git.Entry) error {
	pom, err := maven.PomFromEntry(entry)
	if err != nil {
		return err
	}
	pomInfo, err := pom.FormatInfo()
	if err != nil {
		return err
	}
	parameterizedPomInfo := j.Parameterize(pomInfo)
	for _, p := range parameterizedPomInfo {
		fmt.Println(p)
		if jenkins.IsParameterized(p) {
			missing := strings.Join(jenkins.Parameters(p), ", ")
			found := strings.Join(j.ParameterNames(), ", ")
			indent := ""
			if Flags.Verbose {
				indent = "\t\t"
			}
			_, _ = fmt.Fprintf(os.Stderr,
				"%vjob %v missing parameters pom %v: %v (found: %v)\n",
				indent, j.Name(), pom.Path(), missing, found,
			)
		}
	}
	return nil
}
