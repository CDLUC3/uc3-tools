package jenkins

import (
	"fmt"
	"net/url"
)

// ------------------------------------------------------------
// Job

type Job interface {
	Name() string
	LastSuccess() (Build, error)
	Parameters() []Parameter
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
					params = append(params, &p)
				}
			}
		}
		j.parameters = params
	}
	return j.parameters
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

