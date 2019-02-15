package storagenodes

import "github.com/dmolesUC3/uc3-system-info/internal/output"

type Node struct {
	NodeNumber int64
	Service *CloudService
	Container string
}

func (n *Node) Print(format output.Format) {

}