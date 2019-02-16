package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

type MrtConf struct {
	Root     string
	nodeSets []*NodeSet
}

func NewMrtConf(rootPath string) (*MrtConf, error) {
	conf := MrtConf{Root: rootPath}
	nodePropsDir := conf.nodePropsDir()
	if _, err := os.Stat(nodePropsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("nodes directory %#v does not exist; root %#v does not appear to be an mrt-conf-prv clone", nodePropsDir, rootPath)
	}
	return &conf, nil
}

func (mc *MrtConf) GetNode(env, svc string, number int64) (*Node, error) {
	nodeSet, err := mc.GetNodeSet(env, svc)
	if err != nil {
		return nil, err
	}
	if node, ok := nodeSet.NodesByNumber()[number]; ok {
		return node, nil
	}
	return nil, fmt.Errorf("node number %d not found in node properties %v", number, nodeSet.PropsFilename())
}

func (mc *MrtConf) GetNodeSet(env, svc string) (*NodeSet, error) {
	nodeSets, err := mc.NodeSets()
	if err != nil {
		return nil, err
	}

	propsFile := fmt.Sprintf("nodes-%v-%v.properties", env, svc)
	for _, nodeSet := range nodeSets {
		if propsFile == nodeSet.PropsFilename() {
			return nodeSet, nil
		}
	}
	propsPath := filepath.Join(mc.nodePropsDir(), propsFile)
	return nil, fmt.Errorf("no matching props file %v found for environment %#v, service %#v", propsPath, env, svc)
}


func (mc *MrtConf) NodeSets() ([]*NodeSet, error) {
	if mc.nodeSets == nil {
		nodeSet, err := LoadAllNodes(mc.nodePropsDir())
		if err != nil {
			return nil, err
		}
		mc.nodeSets = nodeSet
	}
	return mc.nodeSets, nil
}

func (mc *MrtConf) s3Conf() string {
	return filepath.Join(mc.Root, "s3-conf")
}

func (mc *MrtConf) s3Resources() string {
	return filepath.Join(mc.s3Conf(), "src", "main", "resources")
}

func (mc *MrtConf) nodePropsDir() string {
	return filepath.Join(mc.s3Resources(), "nodes")
}

