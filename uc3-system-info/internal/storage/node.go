package storage

import (
	"fmt"
	"github.com/CDLUC3/uc3-tools/uc3-system-info/internal/output"
)

type Node struct {
	NodeNumber int64
	Service *CloudService
	Container string
}

func (n *Node) Sprint(format output.Format) string {
	var svcName string
	container := n.Container
	if n.Service == nil {
		svcName = ""
	} else {
		svcName = n.Service.Name
		container = n.Service.ServiceType.ContainerGeneric(container)
	}
	str, _ := format.Sprint(n.NodeNumber, svcName, container)
	return str
}

func (n *Node) ContainerFor(ark string) (string, error) {
	if n.Service == nil {
		return "", fmt.Errorf("no cloud service defined for node %v", n.Description())
	}
	serviceType := n.Service.ServiceType
	return serviceType.ContainerFor(n.Container, ark), nil
}

func (n *Node) KeyFor(ark string, version int, filepath string) string {
	return fmt.Sprintf("%v|%d|%v", ark, version, filepath)
}

func (n *Node) CLIExampleFile(ark string, version int, filepath string) (string, error) {
	service := n.Service
	if service == nil {
		return "", fmt.Errorf("no cloud service defined for node %v", n.Description())
	}
	container, err := n.ContainerFor(ark)
	if err != nil {
		return "", err
	}
	key := n.KeyFor(ark, version, filepath)

	example, err := service.ServiceType.CLIExampleFile(service.Credentials, service.Endpoint, container, key)
	if err != nil {
		return err.Error(), nil
	}
	return example, nil
}

func (n *Node) CLIExampleObject(ark string) (string, error) {
	service := n.Service
	if service == nil {
		return "", fmt.Errorf("no cloud service defined for node %v", n.Description())
	}
	container, err := n.ContainerFor(ark)
	if err != nil {
		return "", err
	}
	example, err := service.ServiceType.CLIExampleObject(service.Credentials, service.Endpoint, container, ark)
	if err != nil {
		return err.Error(), nil
	}
	return example, nil
}

func (n *Node) Description() string {
	return n.Sprint(output.CSV)
}
