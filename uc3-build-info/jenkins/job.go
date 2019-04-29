package jenkins

import (
	"fmt"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/git"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/maven"
	"net/url"
	"sort"
	"strings"
)

// ------------------------------------------------------------
// Job

type Job interface {
	Name() string
	LastSuccess() (Build, error)
	Config() (Config, error)
	ConfigUrl() *url.URL
	APIUrl() *url.URL
	SCMUrl() (string, error)
	Parameters() []Parameter
	Repository() (git.Repository, error)
	POMs() (poms []maven.Pom, errors []error)
}

type JobsByName []Job

func (s JobsByName) Len() int           { return len(s) }
func (s JobsByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s JobsByName) Less(i, j int) bool { return strings.Compare(s[i].Name(), s[j].Name()) < 0 }
func (s JobsByName) Eq(i, j int) bool   { return s[i] == s[j] }
func (s JobsByName) Copy(j, i int)      { s[j] = s[i] }

func (s JobsByName) String() string {
	info := make([]string, len(s))
	for i, dep := range s {
		info[i] = dep.Name()
	}
	return strings.Join(info, ", ")
}

func (s JobsByName) Named(name string) Job {
	for _, j := range s {
		if j.Name() == name {
			return j
		}
	}
	return nil
}

// ------------------------------------------------------------
// Unexported symbols

type job struct {
	JobName             string `json:"name"`
	URL                 string
	LastSuccessfulBuild *build

	Actions []action

	parameters []Parameter
	apiUrl     *url.URL
	configUrl  *url.URL

	config   Config
	repo     git.Repository
	repoPOMs []maven.Pom
}

func (j *job) Name() string {
	return j.JobName
}

func (j *job) String() string {
	return j.Name()
}

func (j *job) LastSuccess() (Build, error) {
	if j.LastSuccessfulBuild == nil {
		if err := j.load(); err != nil {
			return nil, err
		}
		if j.LastSuccessfulBuild == nil {
			return nil, fmt.Errorf("no successful build for job %#v", j.JobName)
		}
	}
	return j.LastSuccessfulBuild, nil
}

func (j *job) Config() (Config, error) {
	if j.config == nil {
		config, err := ConfigFromURL(j.ConfigUrl())
		if err != nil {
			return nil, err
		}
		j.config = config
	}
	return j.config, nil
}

func (j *job) Parameters() []Parameter {
	if j.parameters == nil {
		var params []Parameter
		for _, a := range j.Actions {
			if a.Class == "hudson.model.ParametersDefinitionProperty" {
				for _, p := range a.ParameterDefinitions {
					p1 := p
					sort.Strings(p.Choices_)
					params = append(params, &p1)
				}
			}
		}
		sort.Slice(params, func(i, j int) bool {
			n1, n2 := params[i].Name(), params[j].Name()
			if strings.Compare(n1, n2) < 0 {
				return true
			}
			return false
		})
		j.parameters = params
	}
	return j.parameters
}

func (j *job) Repository() (git.Repository, error) {
	if j.repo == nil {
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
		j.repo = repo
	}
	return j.repo, nil
}

func (j *job) POMs() ([]maven.Pom, []error) {
	var errors []error
	if j.repoPOMs == nil {
		repo, err := j.Repository()
		if err != nil {
			return nil, []error{err}
		}
		config, err := j.Config()
		if err != nil {
			return nil, []error{err}
		}

		pattern := "pom.xml$"
		buildRoot := config.BuildRoot()
		if buildRoot == "" {
		} else {
			pattern = fmt.Sprintf("^%v/.*%v", buildRoot, pattern)
		}

		entries, errs := repo.Find(pattern, git.Blob)
		errors = append(errors, errs...)

		var poms []maven.Pom
		for _, entry := range entries {
			pom, err := maven.PomFromEntry(entry)
			if err != nil {
				errors = append(errors, err)
			}
			if pom != nil {
				poms = append(poms, pom)
			}
		}
		j.repoPOMs = poms
	}
	return j.repoPOMs, errors
}

func (j *job) ConfigUrl() *url.URL {
	if j.configUrl == nil {
		u, err := url.Parse(j.URL)
		if err != nil {
			panic(err)
		}
		j.configUrl = toConfigUrl(u)
	}
	return j.configUrl
}

func (j *job) SCMUrl() (string, error) {
	b, err := j.LastSuccess()
	if err != nil {
		return "", err
	}
	scmUrl, err := b.SCMUrl()
	if err != nil {
		return "", err
	}
	return scmUrl, nil
}

func (j *job) APIUrl() *url.URL {
	if j.apiUrl == nil {
		u, err := url.Parse(j.URL)
		if err != nil {
			panic(err)
		}
		j.apiUrl = toApiUrl(u)
	}
	return j.apiUrl
}

func (j *job) load() error {
	return unmarshal(j.APIUrl(), j)
}

// ------------------------------------------------------------
// Helper types

type action struct {
	Class                string `json:"_class"`
	ParameterDefinitions []parameterDefinition
}
