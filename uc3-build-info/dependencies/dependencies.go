package dependencies

import (
	"fmt"
	. "github.com/CDLUC3/uc3-tools/mrt-build-info/jenkins"
)

func findJob(name string, jobs []Job) (Job, error) {
	for _, j := range jobs {
		if j.Name() == name {
			return j, nil
		}
	}
	return nil, fmt.Errorf("no such job: %#v", name)
}

