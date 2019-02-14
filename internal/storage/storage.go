package storage

import (
	props "github.com/magiconair/properties"
)



func NodesFromFile(nodePropsPath string) error {
	props, err := props.LoadFile(nodePropsPath, props.ISO_8859_1)
	if err != nil {
		return err
	}
	keys := props.Keys()
	for _, k := range keys {

	}
}