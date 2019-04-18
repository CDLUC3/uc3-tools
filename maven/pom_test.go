package maven

import (
	"github.com/dmolesUC3/mrt-build-info/git"
	"github.com/dmolesUC3/mrt-build-info/shared"
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