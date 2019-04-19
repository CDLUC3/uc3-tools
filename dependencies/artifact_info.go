package dependencies

import (
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	"github.com/dmolesUC3/mrt-build-info/maven"
)

type ArtifactInfo interface {
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
