package maven

import (
	"github.com/beevik/etree"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/git"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/shared"
	. "gopkg.in/check.v1"
)

var _ = Suite(&PomSuite{})

type PomSuite struct {
}

func (s *PomSuite) TestEquality(c *C) {
	r, err := git.GetRepository("CDLUC3", "mrt-store", git.SHA1("af174ac555758a1c639a7a3da39e022d9fdbf3a6"))
	c.Assert(err, IsNil)

	entrySHA1 := git.SHA1("d3ae87b904091324ea42af37581d92ee67415143")
	entryUrl := shared.UrlMustParse("https://api.github.com/repos/CDLUC3/mrt-store/git/blobs/d3ae87b904091324ea42af37581d92ee67415143")
	e := r.GetEntry("store-src/pom.xml", entrySHA1, git.Blob, 4775, entryUrl)

	p1, err := PomFromEntry(e)
	c.Assert(err, IsNil)

	p2, err := PomFromEntry(e)
	c.Assert(err, IsNil)

	pp1 := p1.(*pom)
	pp2 := p2.(*pom)

	c.Assert(pp1 == pp2, Equals, true)
}

func (s *PomSuite) TestArtifact(c *C) {
	r, err := git.GetRepository("CDLUC3", "mrt-store", git.SHA1("af174ac555758a1c639a7a3da39e022d9fdbf3a6"))
	c.Assert(err, IsNil)

	entrySHA1 := git.SHA1("d3ae87b904091324ea42af37581d92ee67415143")
	entryUrl := shared.UrlMustParse("https://api.github.com/repos/CDLUC3/mrt-store/git/blobs/d3ae87b904091324ea42af37581d92ee67415143")
	e := r.GetEntry("store-src/pom.xml", entrySHA1, git.Blob, 4775, entryUrl)

	doc := etree.NewDocument()
	file := "testdata/pom.xml"
	err = doc.ReadFromFile(file)
	c.Assert(err, IsNil) // just to be sure

	var pom Pom = &pom{Entry: e, doc: doc}

	artifact, err := pom.Artifact()
	c.Assert(err, IsNil)

	c.Check(artifact.GroupId(), Equals, "org.cdlib.mrt")
	c.Check(artifact.ArtifactId(), Equals, "mrt-storepub-src")
	c.Check(artifact.Packaging(), Equals, "jar")
	c.Check(artifact.Version(), Equals, "1.0-SNAPSHOT")
}

func (s *PomSuite) TestDependencies(c *C) {
	r, err := git.GetRepository("CDLUC3", "mrt-store", git.SHA1("af174ac555758a1c639a7a3da39e022d9fdbf3a6"))
	c.Assert(err, IsNil)

	entrySHA1 := git.SHA1("d3ae87b904091324ea42af37581d92ee67415143")
	entryUrl := shared.UrlMustParse("https://api.github.com/repos/CDLUC3/mrt-store/git/blobs/d3ae87b904091324ea42af37581d92ee67415143")
	e := r.GetEntry("store-src/pom.xml", entrySHA1, git.Blob, 4775, entryUrl)

	doc := etree.NewDocument()
	file := "testdata/pom.xml"
	err = doc.ReadFromFile(file)
	c.Assert(err, IsNil) // just to be sure

	var pom Pom = &pom{Entry: e, doc: doc}

	deps, errs := pom.Dependencies()
	c.Assert(len(errs), Equals, 0) // just to be sure

	expected := []artifact{
		{groupId: "javax.mail", artifactId: "mail", packaging: "jar", version: "1.4.1"},
		{groupId: "javax.servlet", artifactId: "servlet-api", packaging: "jar", version: "2.5"},
		{groupId: "jaxen", artifactId: "jaxen", packaging: "jar", version: "1.1.1"},
		{groupId: "junit", artifactId: "junit", packaging: "jar", version: "4.5"},
		{groupId: "net.sf", artifactId: "jargs", packaging: "jar", version: "1.0"},
		{groupId: "org.cdlib.mrt", artifactId: "mrt-confs3", packaging: "jar", version: "1.0-SNAPSHOT"},
		{groupId: "org.cdlib.mrt", artifactId: "mrt-core", packaging: "jar", version: "2.0-SNAPSHOT"},
		{groupId: "org.cdlib.mrt", artifactId: "mrt-jena", packaging: "jar", version: "2.0-SNAPSHOT"},
		{groupId: "org.cdlib.mrt", artifactId: "mrt-s3srcpub", packaging: "jar", version: "1.0-SNAPSHOT"},
		{groupId: "org.glassfish.jersey.containers", artifactId: "jersey-container-servlet", packaging: "jar", version: "2.25.1"},
		{groupId: "org.glassfish.jersey.core", artifactId: "jersey-client", packaging: "jar", version: "2.25.1"},
		{groupId: "org.glassfish.jersey.media", artifactId: "jersey-media-multipart", packaging: "jar", version: "2.25.1"},
	}
	c.Assert(len(deps), Equals, len(expected))

	for i, a := range deps {
		a1 := a.(*artifact)
		c.Check(*a1, Equals, expected[i])
	}

}

