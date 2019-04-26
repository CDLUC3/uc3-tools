package jenkins

// ------------------------------------------------------------
// Node

type Node interface {
	Jobs() JobsByName
}

// ------------------------------------------------------------
// Unexported symbols

type node struct {
	AllJobs []job `json:"jobs"`

	jobs []Job
}

func (n *node) Jobs() JobsByName {
	if len(n.jobs) != len(n.AllJobs) {
		jobs := make([]Job, len(n.AllJobs))
			for i, j := range n.AllJobs {
			jCopy := j // iteration variables are reused
			jobs[i] = &jCopy
		}
		n.jobs = jobs
	}
	return JobsByName(n.jobs)
}
