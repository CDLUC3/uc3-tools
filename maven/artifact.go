package maven

import (
	"fmt"
	"github.com/beevik/etree"
)

type Artifact interface {
	fmt.Stringer
	GroupId() string
	ArtifactId() string
	Packaging() string
	Version() string
}

func RootArtifact(doc *etree.Document) (Artifact, error) {
	elem := doc.FindElement("/project")
	if elem == nil {
		return nil, fmt.Errorf("<project> not found")
	}
	return artifactFrom(elem)
}

func Dependencies(doc *etree.Document) ([]Artifact, error) {
	var artifacts []Artifact
	for _, elem := range doc.FindElements("/project/dependencies/dependency") {
		a, err := artifactFrom(elem)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, a)
	}
	return artifacts, nil
}

type artifact struct {
	groupId    string
	artifactId string
	packaging  string
	version    string
}

func (a *artifact) String() string {
	return fmt.Sprintf("%v:%v:%v (%v)", a.GroupId(), a.ArtifactId(), a.Version(), a.Packaging())
}

func artifactFrom(elem *etree.Element) (*artifact, error) {
	fields := []string{"groupId", "artifactId", "packaging", "version"}
	values := map[string]string{}
	for _, f := range fields {
		v := elem.FindElement(f)
		if v == nil {
			if f == "packaging" {
				values[f] = "jar" // treat as default
				continue
			}
			return nil, fmt.Errorf("<%s> not found", f)
		}
		values[f] = v.Text()
	}
	a := artifact{
		groupId: values["groupId"],
		artifactId: values["artifactId"],
		packaging: values["packaging"],
		version: values["version"],
	}
	return &a, nil
}



func (a *artifact) GroupId() string {
	return a.groupId
}

func (a *artifact) ArtifactId() string {
	return a.artifactId
}

func (a *artifact) Packaging() string {
	return a.packaging
}

func (a *artifact) Version() string {
	return a.version
}


