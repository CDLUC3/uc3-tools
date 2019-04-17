package maven

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/dmolesUC3/mrt-build-info/git"
	"net/url"
)

type Pom interface {
	fmt.Stringer
	Artifact() (Artifact, error)
	Path() string
	Repository() git.Repository
	URL() *url.URL
	// Deprecated
	FormatInfo() (string, error)
}

func PomFromEntry(entry git.Entry) (Pom, error) {
	if !isPom(entry) {
		return nil, fmt.Errorf("entry %#v does not appear to be a Maven POM", entry.Path())
	}
	return &pom{Entry: entry}, nil
}

type pom struct {
	git.Entry
	doc *etree.Document
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

// TODO:
//   - Get groupId etc. from parent POMs
func (p *pom) Artifact() (Artifact, error) {
	doc, err := p.document()
	if err != nil {
		return nil, err
	}
	return RootArtifact(doc, p.String())
}