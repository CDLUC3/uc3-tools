package jenkins

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// ------------------------------------------------------------
// Job

type Job interface {
	Name() string
	LastSuccess() (Build, error)
	Parameters() []Parameter
	ParameterNames() []string
	Parameterize(str string) []string
}

// ------------------------------------------------------------
// Unexported symbols

type job struct {
	JobName             string `json:"name"`
	URL                 string
	LastSuccessfulBuild *build

	Actions []action

	parameters []Parameter
	apiUrl *url.URL
}

func (j *job) Name() string {
	return j.JobName
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

func (j *job) Parameters() []Parameter {
	if j.parameters == nil {
		var params []Parameter
		for _, a := range j.Actions {
			if a.Class == "hudson.model.ParametersDefinitionProperty" {
				for _, p := range a.ParameterDefinitions {
					sort.Strings(p.Choices_)
					params = append(params, &p)
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

func (j *job) ParameterNames() []string {
	params := j.Parameters()
	paramNames := make([]string, len(params))
	for i, param := range params {
		paramNames[i] = param.Name()
	}
	return paramNames
}

func (j *job) Parameterize(str string) []string {
	parameterized := []string{str}
	for _, param := range j.Parameters() {
		current := parameterized
		var next []string
		for _, c := range current {
			next = append(next, param.Parameterize(c)...)
		}
		parameterized = next
	}
	return parameterized
}

func (j *job) load() error {
	if j.apiUrl == nil {
		u, err := url.Parse(j.URL)
		if err != nil {
			panic(err)
		}
		j.apiUrl = toApiUrl(u)
	}
	return unmarshal(j.apiUrl, j)
}

type action struct {
	Class string `json:"_class"`
	ParameterDefinitions []parameterDefinition
}

