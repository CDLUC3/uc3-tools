package jenkins

// TODO: use an interface; hide unmarshalling
type Node struct {
	Jobs []Job
}
