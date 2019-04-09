package jenkins

import (
	"net/url"
)

type Build struct {
	Number int
	URL string
	Actions []BuildAction
	MavenArtifacts MavenArtifacts

	parsedUrl *url.URL
}

func (b *Build) ApiUrl() *url.URL {
	return toApiUrl(b.url())
}

type BuildAction struct {
	Class string `json:"_class"`
	RemoteURLs []string
	LastBuiltRevision *Revision
}

type Revision struct {
	SHA1 string
	Branches []Branch `json:"branch"`
}

type Branch struct {
	SHA1 string
	Name string
}

type MavenArtifacts struct {
	ModuleRecords []ModuleRecord
}

type ModuleRecord struct {
	MainArtifact *Artifact
	POMArtifact *Artifact
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
