package storage

import (
	"path/filepath"
)

type MrtConf struct {
	Root string
}

func (mc *MrtConf) s3Conf() string {
	return filepath.Join(mc.Root, "s3-conf")
}

func (mc *MrtConf) s3Resources() string {
	return filepath.Join(mc.s3Conf(), "src", "main", "resources")
}

func (mc *MrtConf) nodes() string {
	return filepath.Join(mc.s3Resources(), "nodes")
}
