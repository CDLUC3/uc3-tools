package storage

// NodeIO is "a paramertized definition for all cloud IO members"
// which means ... a set of related access nodes? or something?
type NodeIO struct {
	Nodes         []AccessNode
	NodesByNumber map[int]AccessNode
}

