package dependencies

import (
	"fmt"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/jenkins"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/maven"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/shared"
)

type ArtifactInfo interface {
	fmt.Stringer
	Job() jenkins.Job
	Pom() maven.Pom
	Artifact() maven.Artifact
}

type artifactInfo struct {
	job      jenkins.Job
	pom      maven.Pom
	artifact maven.Artifact
}

func InfoFor(job jenkins.Job, pom maven.Pom, artifact maven.Artifact) ArtifactInfo {
	return &artifactInfo{job: job, pom: pom, artifact: artifact}
}

func (a *artifactInfo) String() string {
	if !shared.Flags.Verbose {
		return a.artifact.String()
	}
	return fmt.Sprintf("%v (%v, %v)", a.artifact, a.pom, a.job)
}

func (a *artifactInfo) Job() jenkins.Job {
	if a == nil {
		return nil
	}
	return a.job
}

func (a *artifactInfo) Pom() maven.Pom {
	if a == nil {
		return nil
	}
	return a.pom
}

func (a *artifactInfo) Artifact() maven.Artifact {
	if a == nil {
		return nil
	}
	return a.artifact
}
