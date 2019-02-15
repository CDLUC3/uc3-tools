package storage

type ByNodeNumber []*Node

func (nodes ByNodeNumber) Len() int {
	return len(nodes)
}

func (nodes ByNodeNumber) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}

func (nodes ByNodeNumber) Less(i, j int) bool {
	return nodes.compare(nodes[i], nodes[j]) < 0
}

func (nodes ByNodeNumber) compare(n1, n2 *Node) int {
	if n1 == n2 {
		return 0
	}
	if n1 == nil {
		if n2 == nil {
			return 0
		}
		return 1
	}
	if n2 == nil {
		return -1
	}
	if n1.NodeNumber < n2.NodeNumber {
		return -1
	}
	if n2.NodeNumber < n1.NodeNumber {
		return 1
	}
	return 0
}