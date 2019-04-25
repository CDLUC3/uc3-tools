package jenkins

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	inTest = true
	TestingT(t)
}