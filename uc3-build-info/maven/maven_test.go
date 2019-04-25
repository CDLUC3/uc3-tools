package maven

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	inTest = true
	TestingT(t)
}
