package maven

import (
	"github.com/CDLUC3/uc3-tools/uc3-build-info/git"
	"strings"
)

var POMURLs = false

var inTest bool = false

func isPom(entry git.Entry) bool {
	return strings.HasSuffix(entry.Path(), "pom.xml")
}

