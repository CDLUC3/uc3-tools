package jenkins

import (
	"net/url"
)

type Job struct {
	Name                string
	URL                 string
	LastSuccessfulBuild *Build

	parsedUrl *url.URL
}

func (j *Job) ApiUrl() *url.URL {
	return toApiUrl(j.url())
}

// ------------------------------------------------------------
// Unexported symbols

func (j *Job) url() *url.URL {
	if j.parsedUrl == nil {
		u, err := url.Parse(j.URL)
		if err != nil {
			panic(err)
		}
		j.parsedUrl = u
	}
	return j.parsedUrl
}

func (j *Job) load() error {
	err := unmarshal(j.ApiUrl(), j)
	if err != nil {
		build := j.LastSuccessfulBuild
		if build != nil {
			err = build.load()
		}
	}
	return err
}