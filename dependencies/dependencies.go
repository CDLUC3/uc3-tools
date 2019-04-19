package dependencies

import (
	"fmt"
	. "github.com/dmolesUC3/mrt-build-info/jenkins"
	. "github.com/dmolesUC3/mrt-build-info/maven"
	. "github.com/dmolesUC3/mrt-build-info/shared"
	"os"
)

func findJob(name string, jobs []Job) (Job, error) {
	for _, j := range jobs {
		if j.Name() == name {
			return j, nil
		}
	}
	return nil, fmt.Errorf("no such job: %#v", name)
}

func mapPomsToJobs(jobs []Job) (map[Pom]Job, []error) {
	//noinspection GoUnhandledErrorResult
	if Flags.Verbose {
		fmt.Fprintf(os.Stderr, "Retrieving POMs for %d jobs...", len(jobs))
		defer func() { fmt.Fprintln(os.Stderr, "Done.") }()
	}

	var errors []error
	jobsByPom := map[Pom]Job{}
	for _, j := range jobs {
		//noinspection GoUnhandledErrorResult
		if Flags.Verbose {
			fmt.Fprint(os.Stderr, ".")
		}
		poms, errs := j.POMs()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		if len(poms) == 0 {
			errors = append(errors, fmt.Errorf("no POMs found for job %v", j.Name()))
		}
		for _, p := range poms {
			jobsByPom[p] = j
		}
	}
	return jobsByPom, errors
}
