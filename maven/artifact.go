package maven

import (
	"fmt"
	"github.com/beevik/etree"
	. "github.com/dmolesUC3/mrt-build-info/shared"
	"strings"
)

type Artifact interface {
	fmt.Stringer
	GroupId() string
	ArtifactId() string
	Packaging() string
	Version() string
}

type ArtifactsByString []Artifact

func (a ArtifactsByString) Len() int           { return len(a) }
func (a ArtifactsByString) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ArtifactsByString) Less(i, j int) bool { return strings.Compare(a[i].String(), a[j].String()) < 0 }

func (a ArtifactsByString) String() string {
	info := make([]string, len(a))
	for i, dep := range a {
		info[i] = dep.String()
	}
	return strings.Join(info, ", ")
}

var artifactCache = map[artifact]*artifact{}

func GetArtifact(groupId string, artifactId string, packaging string, version string) Artifact {
	arec := artifact{groupId: groupId, artifactId: artifactId, packaging: packaging, version: version}
	if aptr, ok := artifactCache[arec]; ok {
		return aptr
	}
	aptr := &arec
	artifactCache[arec] = aptr
	return aptr
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
func artifactFrom(elem *etree.Element) (Artifact, error) {
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
	return GetArtifact(values["groupId"], values["artifactId"], values["packaging"], values["version"]), nil
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
