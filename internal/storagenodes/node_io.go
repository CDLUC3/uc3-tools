package storagenodes

import (
	"fmt"
	props "github.com/magiconair/properties"
	"path/filepath"
	"strconv"
	"strings"
)

type NodeSet struct {
	propsDir string
	services map[string]*CloudService
	nodes    map[int64]*Node
}

func LoadNodes(propsPath string) (*NodeSet, error) {
	nodeProps, err := props.LoadFile(propsPath, props.ISO_8859_1)
	if err != nil {
		return nil, err
	}

	propsDir := filepath.Clean(filepath.Dir(propsPath))
	ns := NodeSet{propsDir: propsDir}
	keys := nodeProps.Keys()
	for _, k := range keys {
		if v, ok := nodeProps.Get(k); ok {
			node, err := ns.loadNode(v)
			if err != nil {
				return nil, err
			}
			nodeNum := node.NodeNumber
			if _, exists := ns.nodes[nodeNum]; exists {
				return nil, fmt.Errorf("duplicate node number: %d", nodeNum)
			}
			ns.nodes[nodeNum] = node
		}
	}
	return &ns, nil
}

func Concat(ns1, ns2 *NodeSet) (*NodeSet, error) {
	if ns1.propsDir != ns2.propsDir {
		return nil, fmt.Errorf("can't join nodesets from different directories: %v, %v", ns1.propsDir, ns2.propsDir)
	}
	var ns3 NodeSet
	for name, svc1 := range ns1.services {
		if svc2, exists := ns2.services[name]; exists {
			if *svc2 != *svc1 {
				return nil, fmt.Errorf("incompatible service definitions for %v", name)
			}
		}
		svc3 := *svc1
		ns3.services[name] = &svc3
	}
	for name, svc2 := range ns1.services {
		if _, exists := ns3.services[name]; !exists {
			svc3 := *svc2
			ns3.services[name] = &svc3
		}
	}
	for num, node1 := range ns1.nodes {
		if node2, exists := ns2.nodes[num]; exists {
			if *node2 != *node1 {
				return nil, fmt.Errorf("incompatible definitions for node %d", num)
			}
			node3 := *node1
			svc3, err := ns3.serviceFor(node1.Service.Name)
			if err != nil {
				return nil, err
			}
			node3.Service = svc3
			ns3.nodes[num] = &node3
		}
	}
	for num, node2 := range ns2.nodes {
		if _, exists := ns3.nodes[num]; !exists {
			node3 := *node2
			svc3, err := ns3.serviceFor(node2.Service.Name)
			if err != nil {
				return nil, err
			}
			node3.Service = svc3
			ns3.nodes[num] = &node3
		}
	}
	return &ns3, nil
}

// ------------------------------------------------------------
// Unexported functions

func (ns *NodeSet) loadNode(nodeLine string) (*Node, error) {
	fields := strings.Split(nodeLine, "|")
	if len(fields) < 2 {
		return nil, fmt.Errorf("not enough fields: %v=%v", k, nodeLine)
	}
	nodeNum, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return nil, err
	}

	svcName := fields[1]
	svc, err := ns.serviceFor(svcName)
	if err != nil {
		return nil, err
	}

	var container string
	if len(fields) > 2 {
		container = fields[2]
	}

	return &Node{nodeNum, svc, container}, nil
}

func (ns *NodeSet) serviceFor(svcName string) (*CloudService, error) {
	if cs, ok := ns.services[svcName]; ok {
		return cs, nil
	}
	svcPropsPath := filepath.Join(ns.propsDir, svcName + ".properties")
	cs, err := LoadCloudService(svcName, svcPropsPath)
	if err != nil {
		return nil, err
	}
	ns.services[svcName] = cs
	return cs, nil
}
