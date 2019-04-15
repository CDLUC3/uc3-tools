package maven

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/dmolesUC3/mrt-build-info/git"
)

type Pom interface {
	Path() string
	Repository() git.Repository
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

func (p *pom) artifact() (Artifact, error) {
	doc, err := p.document()
	if err != nil {
		return nil, err
	}
	return RootArtifact(doc)
}