package jenkins

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/git"
	"github.com/dmolesUC3/mrt-build-info/maven"
	"net/url"
	"regexp"
)

// ------------------------------------------------------------
// Build

type Build interface {
	BuildNumber() int
	SCMUrl() (string, error)
	SHA1() (git.SHA1, error)
	Artifacts() ([]maven.Artifact, error)
	Commit() (owner, repoName string, sha1 git.SHA1, err error)
}

// ------------------------------------------------------------
// Unexported symbols

const buildDataClass = "hudson.plugins.git.util.BuildData"

var repoRe = regexp.MustCompile("/(.+)/([^./]+)(?:\\.git)?$")

type build struct {
	Number          int
	URL             string
	Actions         []buildAction
	MavenArtifacts  *mavenArtifacts
	FullDisplayName string

	apiUrl          *url.URL
	artifacts       []maven.Artifact
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

func (b *build) Commit() (owner, repoName string, sha1 git.SHA1, err error) {
	scm, err := b.SCMUrl()
	if err == nil {
		u, err := url.Parse(scm)
		if err == nil {
			if repoRe.MatchString(u.Path) {
				matches := repoRe.FindStringSubmatch(u.Path)
				owner = matches[1]
				repoName = matches[2]
				sha1, err = b.SHA1()
			} else {
				err = fmt.Errorf("SCM URL %#v for %v does not appear to be a Git URL", u, b.FullDisplayName)
			}
		}
	}
	return
}

func (b *build) SHA1() (git.SHA1, error) {
	bd, err := b.buildData()
	if err != nil {
		return "", err
	}
	rev := bd.LastBuiltRevision
	if rev == nil {
		return "", fmt.Errorf("can't find revision for %v", b.FullDisplayName)
	}
	return rev.SHA1, nil
}

func (b *build) Artifacts() ([]maven.Artifact, error) {
	if b.artifacts == nil {
		if b.MavenArtifacts == nil {
			if err := b.load(); err != nil {
				return nil, err
			}
		}
		if b.MavenArtifacts == nil {
			b.artifacts = make([]maven.Artifact, 0)
		} else {
			moduleRecords := b.MavenArtifacts.ModuleRecords
			artifacts := make([]maven.Artifact, len(moduleRecords))
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
	SHA1     git.SHA1
	Branches []branch `json:"branch"`
}

type branch struct {
	SHA1 string
	Name string
}

// ------------------------------------------------------------
// Maven artifact information

type artifact struct {
	GroupId_    string `json:"groupId"`
	ArtifactId_ string `json:"artifactId"`
	Type        string `json:"type"`
	Version_    string `json:"version"`
}

func (a *artifact) String() string {
	return maven.ArtifactToString(a)
}

func (a *artifact) GroupId() string {
	return a.GroupId_
}

func (a *artifact) ArtifactId() string {
	return a.ArtifactId_
}

func (a *artifact) Packaging() string {
	return a.Type
}

func (a *artifact) Version() string {
	return a.Version_
}

type mavenArtifacts struct {
	ModuleRecords []moduleRecord
}

type moduleRecord struct {
	MainArtifact *artifact
	POMArtifact  *artifact
}
