package maven

import (
	"github.com/beevik/etree"
	. "gopkg.in/check.v1"
)

var _ = Suite(&ArtifactSuite{})

type ArtifactSuite struct {
}

func (s *ArtifactSuite) TestRootArtifact(c *C) {
	doc := etree.NewDocument()
	err := doc.ReadFromFile("testdata/pom.xml")
	c.Assert(err, IsNil) // just to be sure

	artifact, err := RootArtifact(doc)
	c.Assert(err, IsNil)

	c.Check(artifact.GroupId(), Equals, "org.cdlib.mrt")
	c.Check(artifact.ArtifactId(), Equals, "mrt-storepub-src")
	c.Check(artifact.Packaging(), Equals, "jar")
	c.Check(artifact.Version(), Equals, "1.0-SNAPSHOT")
}

func (s *ArtifactSuite) TestDependencies(c *C) {
	doc := etree.NewDocument()
	err := doc.ReadFromFile("testdata/pom.xml")
	c.Assert(err, IsNil) // just to be sure

	artifacts, err := Dependencies(doc)
	c.Assert(err, IsNil)

	expected := []artifact{
		{groupId: "org.glassfish.jersey.containers", artifactId: "jersey-container-servlet", packaging: "jar", version: "2.25.1"},
		{groupId: "org.glassfish.jersey.media", artifactId: "jersey-media-multipart", packaging: "jar", version: "2.25.1"},
		{groupId: "org.glassfish.jersey.core", artifactId: "jersey-client", packaging: "jar", version: "2.25.1"},
		{groupId: "org.cdlib.mrt", artifactId: "mrt-core", packaging: "jar", version: "2.0-SNAPSHOT"},
		{groupId: "org.cdlib.mrt", artifactId: "mrt-jena", packaging: "jar", version: "2.0-SNAPSHOT"},
		{groupId: "org.cdlib.mrt", artifactId: "mrt-s3srcpub", packaging: "jar", version: "1.0-SNAPSHOT"},
		{groupId: "org.cdlib.mrt", artifactId: "mrt-confs3", packaging: "jar", version: "1.0-SNAPSHOT"},
		{groupId: "junit", artifactId: "junit", packaging: "jar", version: "4.5"},
		{groupId: "net.sf", artifactId: "jargs", packaging: "jar", version: "1.0"},
		{groupId: "javax.servlet", artifactId: "servlet-api", packaging: "jar", version: "2.5"},
		{groupId: "jaxen", artifactId: "jaxen", packaging: "jar", version: "1.1.1"},
		{groupId: "javax.mail", artifactId: "mail", packaging: "jar", version: "1.4.1"},
	}
	c.Assert(len(artifacts), Equals, len(expected))

	for i, a := range artifacts {
		a1 := a.(*artifact)
		c.Check(*a1, Equals, expected[i])
	}
}
