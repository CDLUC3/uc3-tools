package git

import (
	"fmt"
	"github.com/google/go-github/github"
)

type EntryType uint

const (
	Blob EntryType = iota
	Tree
)

func GetEntryType(entry github.TreeEntry) EntryType {
	switch entry.GetType() {
	case "blob":
		return Blob
	case "tree":
		return Tree
	default:
		return 0
	}
}

func (e EntryType) String() string {
	switch e {
	case Blob:
		return "Blob"
	case Tree:
		return "Tree"
	default:
		return fmt.Sprintf("unknown (%x)", uint(e))
	}
}
