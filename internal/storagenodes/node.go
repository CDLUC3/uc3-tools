package storagenodes

type Node struct {
	NodeNumber int64
	Service *CloudService
	Container string
}