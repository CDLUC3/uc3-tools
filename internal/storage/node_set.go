package storage

import (
	"fmt"
	props "github.com/magiconair/properties"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type NodeSet struct {
	PropsPath     string
	services      map[string]*CloudService
	nodesByNumber map[int64]*Node
	sortedNodes   []*Node
}

func (ns *NodeSet) Nodes() (nodes []*Node) {
	if ns == nil {
		return nodes
	}
	if len(ns.sortedNodes) == 0 {
		var sn []*Node
		for _, v := range ns.NodesByNumber() {
			sn = append(sn, v)
		}
		sort.Sort(ByNodeNumber(sn))
		ns.sortedNodes = sn
	}
	return ns.sortedNodes
}

func (ns *NodeSet) Services() map[string]*CloudService {
	if ns.services == nil {
		ns.services = map[string]*CloudService{}
	}
	return ns.services
}

func (ns *NodeSet) NodesByNumber() map[int64]*Node {
	if ns.nodesByNumber == nil {
		ns.nodesByNumber = map[int64]*Node{}
	}
	return ns.nodesByNumber
}

func (ns *NodeSet) PropsFilename() string {
	return filepath.Base(ns.PropsPath)
}

func LoadAllNodes(propsDir string) ([]*NodeSet, error) {
	propsPaths, err := filepath.Glob(filepath.Join(propsDir, "nodes-*.properties"))
	if err != nil {
		return nil, err
	}
	if len(propsPaths) == 0 {
		return nil, fmt.Errorf("no nodes-*.properties files found in %v", propsDir)
	}

	var nodes []*NodeSet
	for _, propsPath := range propsPaths {
		ns, err := LoadNodes(propsPath)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, ns)
	}
	return nodes, nil
}

func LoadNodes(propsPath string) (*NodeSet, error) {
	nodeProps, err := props.LoadFile(propsPath, props.ISO_8859_1)
	if err != nil {
		return nil, err
	}

	ns := NodeSet{PropsPath: propsPath}
	keys := nodeProps.Keys()
	for _, k := range keys {
		if v, ok := nodeProps.Get(k); ok {
			node, err := ns.loadNode(k, v)
			if err != nil {
				// TODO: verbose logging?
				//_, _ = fmt.Fprintf(os.Stderr, "invalid node definition %v=%v in %v: %v\n", k, v, filepath.Base(propsPath), err.Error())
				continue
			}
			nodeNum := node.NodeNumber
			if _, exists := ns.NodesByNumber()[nodeNum]; exists {
				return nil, fmt.Errorf("duplicate node number: %d", nodeNum)
			}
			ns.NodesByNumber()[nodeNum] = node
		}
	}
	return &ns, nil
}

// ------------------------------------------------------------
// Unexported functions

func (ns *NodeSet) loadNode(nodeSeq, nodeLine string) (*Node, error) {
	fields := strings.Split(nodeLine, "|")
	if len(fields) < 2 {
		return nil, fmt.Errorf("not enough fields: %v=%v", nodeSeq, nodeLine)
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
	services := ns.Services()
	if cs, ok := services[svcName]; ok {
		return cs, nil
	}
	propsDir := filepath.Dir(ns.PropsPath)
	svcPropsPath := filepath.Join(propsDir, svcName + ".properties")
	cs, err := LoadCloudService(svcName, svcPropsPath)
	if err != nil {
		return nil, err
	}
	services[svcName] = cs
	return cs, nil
}
