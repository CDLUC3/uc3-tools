package maven

import (
	"fmt"
	"github.com/beevik/etree"
	. "github.com/dmolesUC3/mrt-build-info/shared"
)

type Artifact interface {
	fmt.Stringer
	GroupId() string
	ArtifactId() string
	Packaging() string
	Version() string
}

// TODO: something smarter than passing source around
func RootArtifact(doc *etree.Document, source string) (Artifact, error) {
	elem := doc.FindElement("/project")
	if elem == nil {
		return nil, fmt.Errorf("<project> not found in " + source)
	}
	return artifactFrom(elem, source)
}

// TODO: something smarter than passing source around
func Dependencies(doc *etree.Document, source string) ([]Artifact, error) {
	var artifacts []Artifact
	for _, elem := range doc.FindElements("/project/dependencies/dependency") {
		a, err := artifactFrom(elem, source)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, a)
	}
	return artifacts, nil
}

func ArtifactToString(a Artifact) string {
	if Flags.Verbose {
		return fmt.Sprintf("%v:%v:%v (%v)", a.GroupId(), a.ArtifactId(), a.Version(), a.Packaging())
	}
	return fmt.Sprintf("%v:%v:%v", a.GroupId(), a.ArtifactId(), a.Version())
}

type artifact struct {
	groupId    string
	artifactId string
	packaging  string
	version    string
}

func (a *artifact) String() string {
	return ArtifactToString(a)
}

// TODO:
//   - for root, get groupId etc. from parent POMs
func artifactFrom(elem *etree.Element, source string) (*artifact, error) {
	fields := []string{"groupId", "artifactId", "packaging", "version"}
	values := map[string]string{}
	for _, f := range fields {
		v := elem.FindElement(f)
		if v == nil {
			if f == "packaging" {
				values[f] = "jar" // treat as default
				continue
			}
			return nil, fmt.Errorf("<%s> not found in %v", f, source)
		}
		values[f] = v.Text()
	}
	a := artifact{
		groupId:    values["groupId"],
		artifactId: values["artifactId"],
		packaging:  values["packaging"],
		version:    values["version"],
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
