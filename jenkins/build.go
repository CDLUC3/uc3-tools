package jenkins

import (
	"fmt"
	"net/url"
	"regexp"
)

// ------------------------------------------------------------
// Build

type Build interface {
	BuildNumber() int
	SCMUrl() (string, error)
	SHA1() (string, error)
	Artifacts() ([]Artifact, error)
	Repo() (string, error)
}

type Artifact interface {
	Group() string
	Artifact() string
	Type() string
	Version() string
	File() string
}

// ------------------------------------------------------------
// Unexported symbols

const buildDataClass = "hudson.plugins.git.util.BuildData"
var repoRe = regexp.MustCompile("/(.+)\\.git")

type build struct {
	Number         int
	URL            string
	Actions        []buildAction
	MavenArtifacts *mavenArtifacts
	FullDisplayName string

	apiUrl    *url.URL
	artifacts []Artifact
	buildDataAction *buildAction
}

func (b *build) BuildNumber() int {
	return b.Number
}

func (b *build) SCMUrl() (string, error) {
	bd, err := b.buildData()
	if err != nil {
		return "", err
	}
	if len(bd.RemoteURLs) != 1 {
		return "", fmt.Errorf("can't determine remote URL for %v: expected 1 remote URL, found %d", b.FullDisplayName, len(bd.RemoteURLs))
	}
	return bd.RemoteURLs[0], nil
}


func (b *build) Repo() (string, error) {
	scm, err := b.SCMUrl()
	if err != nil {
		return "", err
	}
	u, err := url.Parse(scm)
	if err != nil {
		return "", err
	}
	if !repoRe.MatchString(u.Path) {
		return "", fmt.Errorf("SCM URL %#v for %v does not appear to be a Git URL", u, b.FullDisplayName)
	}
	return repoRe.FindStringSubmatch(u.Path)[1], nil
}


func (b *build) SHA1() (string, error) {
	bd, err := b.buildData()
	if err != nil {
		return "", err
	}
	rev := bd.LastBuiltRevision
	if rev == nil {
		return "", fmt.Errorf("can't revision for %v", b.FullDisplayName)
	}
	return rev.SHA1, nil
}

func (b *build) Artifacts() ([]Artifact, error) {
	if b.artifacts == nil {
		if b.MavenArtifacts == nil {
			if err := b.load(); err != nil {
				return nil, err
			}
		}
		if b.MavenArtifacts == nil {
			b.artifacts = make([]Artifact, 0)
		} else {
			moduleRecords := b.MavenArtifacts.ModuleRecords
			artifacts := make([]Artifact, len(moduleRecords))
			for i, r := range moduleRecords {
				artifacts[i] = r.MainArtifact
			}
			b.artifacts = artifacts
		}
	}
	return b.artifacts, nil
}

func (b *build) buildData() (*buildAction, error) {
	if b.buildDataAction == nil {
		if b.Actions == nil {
			if err := b.load(); err != nil {
				return nil, err
			}
		}
		for _, a := range b.Actions {
			if a.Class == buildDataClass {
				b.buildDataAction = &a
				break
			}
		}
		if b.buildDataAction == nil {
			return nil, fmt.Errorf("%v action not found in %v (%v)", buildDataClass, b.FullDisplayName, b.apiUrl)
		}
	}
	return b.buildDataAction, nil
}

func (b *build) load() error {
	if b.apiUrl == nil {
		u, err := url.Parse(b.URL)
		if err != nil {
			panic(err)
		}
		b.apiUrl = toApiUrl(u)
	}
	return unmarshal(b.apiUrl, b)
}

// ------------------------------------------------------------
// SCM information

type buildAction struct {
	Class             string `json:"_class"`
	RemoteURLs        []string
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

type artifact struct {
	GroupId         string
	ArtifactId      string
	ArtifactType    string `json:"type"`
	ArtifactVersion string `json:"version"`
	FileName        string
}

func (a *artifact) Group() string {
	return a.GroupId
}

func (a *artifact) Artifact() string {
	return a.ArtifactId
}

func (a *artifact) Type() string {
	return a.ArtifactType
}

func (a *artifact) Version() string {
	return a.ArtifactVersion
}

func (a *artifact) File() string {
	return a.FileName
}

type mavenArtifacts struct {
	ModuleRecords []moduleRecord
}

type moduleRecord struct {
	MainArtifact *artifact
	POMArtifact  *artifact
}
