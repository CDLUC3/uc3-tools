package maven

import (
	. "gopkg.in/check.v1"
)

var _ = Suite(&ArtifactSuite{})

type ArtifactSuite struct {
}

func (s *ArtifactSuite) TestEquality(c *C) {
	a1 := GetArtifact("org.cdlib.mrt", "mrt-storepub-src", "jar", "1.0-SNAPSHOT")
	a2 := GetArtifact("org.cdlib.mrt", "mrt-storepub-src", "jar", "1.0-SNAPSHOT")

	c.Assert(a1 == a2, Equals, true)

	aa1 := a1.(*artifact)
	aa2 := a2.(*artifact)

	c.Assert(aa1 == aa2, Equals, true)
	c.Assert(*aa1 == *aa2, Equals, true)
}
