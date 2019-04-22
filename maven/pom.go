package maven

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/dmolesUC3/mrt-build-info/git"
	"net/url"
	"sort"
	"strings"
)

var pomCache = map[git.Entry]Pom{}

type Pom interface {
	fmt.Stringer
	Artifact() (Artifact, error)
	Dependencies() ([]Artifact, []error)
	Path() string
	Repository() git.Repository
	Location() string
	URL() *url.URL
	BlobURL() *url.URL
	// Deprecated
	FormatInfo() (string, error)
}

type PomsByLocation []Pom
func (p PomsByLocation) Len() int           { return len(p) }
func (p PomsByLocation) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PomsByLocation) Less(i, j int) bool { return strings.Compare(p[i].Location(), p[j].Location()) < 0 }

func (p PomsByLocation) String() string {
	info := make([]string, len(p))
	for i, dep := range p {
		info[i] = dep.Location()
	}
	return strings.Join(info, ", ")
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

	artifact     Artifact
	dependencies []Artifact
}

func (p *pom) String() string {
	return p.Location()
}

func (p *pom) Location() string {
	return fmt.Sprintf("%v/%v", p.Repository(), p.Path())
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
			return nil, p.addLocation(err)
		}
		doc := etree.NewDocument()
		err = doc.ReadFromBytes(data)
		if err != nil {
			return nil, p.addLocation(err)
		}
		p.doc = doc
	}
	return p.doc, nil
}

func (p *pom) Artifact() (Artifact, error) {
	if p.artifact == nil {
		doc, err := p.document()
		if err != nil {
			return nil, err
		}
		elem := doc.FindElement("/project")
		if elem == nil {
			return nil, fmt.Errorf("<project> not found in %v", p.Location())
		}
		artifact, err := artifactFrom(elem)
		if err != nil {
			return nil, p.addLocation(err)
		}
		p.artifact = artifact
	}
	return p.artifact, nil
}

func (p *pom) Dependencies() ([]Artifact, []error) {
	var errors []error
	if p.dependencies == nil {
		doc, err := p.document()
		if err != nil {
			errors = append(errors, err)
			p.dependencies = []Artifact{}
		} else {
			var deps []Artifact
			for _, elem := range doc.FindElements("/project/dependencies/dependency") {
				a, err := artifactFrom(elem)
				if err != nil {
					errors = append(errors, p.addLocation(err))
					continue
				}
				deps = append(deps, a)
			}
			sort.Sort(ArtifactsByString(deps))
			p.dependencies = deps
		}
	}
	return p.dependencies, errors
}

func (p *pom) addLocation(err error) error {
	return fmt.Errorf("%v in %v", err, p.Location())
}

