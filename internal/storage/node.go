package storage

import "github.com/dmolesUC3/uc3-system-info/internal/output"

type Node struct {
	NodeNumber int64
	Service *CloudService
	Container string
}

func (n *Node) Sprint(format output.Format) string {
	var svcName string
	if n.Service == nil {
		svcName = "<nil>"
	} else {
		svcName = n.Service.Name
	}
	str, _ := format.Sprint(n.NodeNumber, svcName, n.Container)
	return str
}