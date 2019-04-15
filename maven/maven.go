package maven

import (
	"github.com/dmolesUC3/mrt-build-info/git"
	"strings"
)

var inTest bool = false

func isPom(entry git.Entry) bool {
	return strings.HasSuffix(entry.Path(), "pom.xml")
}

