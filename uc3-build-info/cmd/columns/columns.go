package columns

import (
	"github.com/CDLUC3/uc3-tools/mrt-build-info/jenkins"
	"github.com/CDLUC3/uc3-tools/mrt-build-info/maven"
	. "github.com/CDLUC3/uc3-tools/mrt-build-info/shared"
)

const ValueUnknown = "(unknown)"

func Job(jobs []jenkins.Job) TableColumn {
	return Jobs(func(row int) jenkins.Job { return jobs[row] }, len(jobs))
}

func Pom(poms []maven.Pom) TableColumn {
	return Poms(func(row int) maven.Pom { return poms[row] }, len(poms))
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

