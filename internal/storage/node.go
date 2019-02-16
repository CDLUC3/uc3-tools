package storage

import (
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
)

type Node struct {
	NodeNumber int64
	Service *CloudService
	Container string
}

func (n *Node) Sprint(format output.Format) string {
	var svcName string
	if n.Service == nil {
		svcName = ""
	} else {
		svcName = n.Service.Name
	}
	str, _ := format.Sprint(n.NodeNumber, svcName, n.Container)
	return str
}

func (n *Node) ContainerFor(ark string) (string, error) {
	if n.Service == nil {
		return "", fmt.Errorf("no cloud service defined for node %v", n.Sprint(output.CSV))
	}
	serviceType := n.Service.ServiceType
	return serviceType.ContainerFor(n.Container, ark), nil
}