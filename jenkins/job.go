package jenkins

import (
	"net/url"
)

// ------------------------------------------------------------
// Job

type Job interface {
	Name() string
	LastSuccess() (Build, error)
}

// ------------------------------------------------------------
// Unexported symbols

type job struct {
	JobName             string `json:"name"`
	URL                 string
	LastSuccessfulBuild *build

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
	}
	return j.LastSuccessfulBuild, nil
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


