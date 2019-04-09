package jenkins

import (
	"net/url"
)

type Build struct {
	actions		// TODO: hide unexported types more effectively
	artifacts	// TODO: hide unexported types more effectively

	Number         int
	URL            string

	parsedUrl *url.URL
}


func (b *Build) ApiUrl() *url.URL {
	return toApiUrl(b.url())
}

type Artifact struct {
	GroupId string
	ArtifactId string
	Md5Sum string
	Type string
	Version string
	CanonicalName string
	FileName string
}

// ------------------------------------------------------------
// Unexported symbols

func (b *Build) url() *url.URL {
	if b.parsedUrl == nil {
		u, err := url.Parse(b.URL)
		if err != nil {
			panic(err)
		}
		b.parsedUrl = u
	}
	return b.parsedUrl
}

func (b *Build) load() error {
	return unmarshal(b.ApiUrl(), b)
}

// ------------------------------------------------------------
// SCM information

type actions struct {
	Actions        []buildAction
}

type buildAction struct {
	Class string `json:"_class"`
	RemoteURLs []string
	LastBuiltRevision *revision
}

type revision struct {
	SHA1     string
	Branches []branch `json:"branch"`
}

type branch struct {
	SHA1 string
	Name string
}

// ------------------------------------------------------------
// Maven artifact information

type artifacts struct {
	MavenArtifacts mavenArtifacts
}

type mavenArtifacts struct {
	ModuleRecords []moduleRecord
}

type moduleRecord struct {
	MainArtifact *Artifact
	POMArtifact  *Artifact
}

