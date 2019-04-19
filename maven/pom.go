package maven

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/dmolesUC3/mrt-build-info/git"
	"net/url"
	"sort"
)

var pomCache = map[git.Entry]Pom{}

type Pom interface {
	fmt.Stringer
	Artifact() (Artifact, error)
	Dependencies() ([]Artifact, error)
	Path() string
	Repository() git.Repository
	URL() *url.URL
	BlobURL() *url.URL
	// Deprecated
	FormatInfo() (string, error)
}

func PomFromEntry(entry git.Entry) (Pom, error) {
	if !isPom(entry) {
		return nil, fmt.Errorf("entry %#v does not appear to be a Maven POM", entry.Path())
	}
	if p, ok := pomCache[entry]; ok {
		return p, nil
	}
	p := &pom{Entry: entry}
	pomCache[entry] = p
	return p, nil
}

type pom struct {
	git.Entry
	doc *etree.Document

	dependencies []Artifact
}

func (p *pom) String() string {
	return fmt.Sprintf("%v (%v)", p.Path(), p.Repository())
}

// Deprecated
func (p *pom) FormatInfo() (string, error) {
	artifact, err := p.Artifact()
	if err != nil {
		return "", err
	}
	pomInfo := fmt.Sprintf("%v\t%v\t%v", artifact.String(), p.Repository(), p.Path())
	if POMURLs {
		pomInfo = fmt.Sprintf("%v\t%v", pomInfo, p.URL())
	}
	return pomInfo, nil
}

func (p *pom) URL() *url.URL {
	return git.WebUrlForEntry(p.Entry)
}

func (p *pom) BlobURL() *url.URL {
	return p.Entry.URL()
}

func (p *pom) document() (*etree.Document, error) {
	if p.doc == nil {
		data, err := p.GetContent()
		if err != nil {
			return nil, err
		}
		doc := etree.NewDocument()
		err = doc.ReadFromBytes(data)
		if err != nil {
			return nil, err
		}
		p.doc = doc
	}
	return p.doc, nil
}

func (p *pom) Artifact() (Artifact, error) {
	doc, err := p.document()
	if err != nil {
		return nil, err
	}
	return RootArtifact(doc, p.String())
}

func (p *pom) Dependencies() ([]Artifact, error) {
	if p.dependencies == nil {
		doc, err := p.document()
		if err != nil {
			return nil, err
		}
		deps, err := Dependencies(doc, p.Path())
		if err != nil {
			return nil, err
		}
		sort.Sort(ArtifactsByString(deps))
		p.dependencies = deps
	}
	return p.dependencies, nil
}
