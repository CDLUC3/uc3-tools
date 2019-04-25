package dependencies

import (
	"fmt"
	. "github.com/dmolesUC3/mrt-build-info/jenkins"
)

func findJob(name string, jobs []Job) (Job, error) {
	for _, j := range jobs {
		if j.Name() == name {
			return j, nil
		}
	}
	return nil, fmt.Errorf("no such job: %#v", name)
}

