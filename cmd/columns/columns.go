package columns

import (
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	"github.com/dmolesUC3/mrt-build-info/maven"
	. "github.com/dmolesUC3/mrt-build-info/shared"
)

const ValueUnknown = "(unknown)"

func Job(jobs []jenkins.Job) TableColumn {
	return Jobs(func(row int) jenkins.Job { return jobs[row] }, len(jobs))
}

func Pom(poms []maven.Pom) TableColumn {
	return Poms(func(row int) maven.Pom { return poms[row] }, len(poms))
}

func Artifacts(artFor func(row int) maven.Artifact, rows int) TableColumn {
	return NewTableColumn("Artifact", rows, func(row int) string {
		art := artFor(row)
		if art == nil {
			return ""
		}
		return art.String()
	})
}

func Jobs(jobFor func(row int) jenkins.Job, rows int) TableColumn {
	return NewTableColumn("Job Name", rows, func(row int) string {
		job := jobFor(row)
		if job == nil {
			return ""
		}
		return job.Name()
	})
}

func Poms(pomFor func(row int) maven.Pom, rows int) TableColumn {
	return NewTableColumn("POM", rows, func(row int) string {
		pom := pomFor(row)
		if pom == nil {
			return ""
		}
		return pom.Path()
	})
}

