package git

import (
	"github.com/dmolesUC3/mrt-build-info/shared"
	. "gopkg.in/check.v1"
)

var _ = Suite(&RepositorySuite{})

type RepositorySuite struct {}

func (s *RepositorySuite) TestEquality(c *C) {
	r1, err := GetRepository("CDLUC3", "mrt-store", SHA1("af174ac555758a1c639a7a3da39e022d9fdbf3a6"))
	c.Assert(err, IsNil)

	r2, err := GetRepository("cdluc3", "mrt-store", SHA1("af174ac555758a1c639a7a3da39e022d9fdbf3a6"))
	c.Assert(err, IsNil)

	c.Assert(r1, Equals, r2)

	rr1 := r1.(*repository)
	rr2 := r2.(*repository)
	c.Assert(rr1 == rr2, Equals, true)
}

func (s *RepositorySuite) TestEntry(c *C) {
	r1, err := GetRepository("CDLUC3", "mrt-store", SHA1("af174ac555758a1c639a7a3da39e022d9fdbf3a6"))
	c.Assert(err, IsNil)

	r2, err := GetRepository("cdluc3", "mrt-store", SHA1("af174ac555758a1c639a7a3da39e022d9fdbf3a6"))
	c.Assert(err, IsNil)

	e1 := r1.GetEntry("store-src/pom.xml", SHA1("d3ae87b904091324ea42af37581d92ee67415143"), Blob, 4775, shared.UrlMustParse("https://api.github.com/repos/CDLUC3/mrt-store/git/blobs/d3ae87b904091324ea42af37581d92ee67415143"))
	e2 := r2.GetEntry("store-src/pom.xml", SHA1("d3ae87b904091324ea42af37581d92ee67415143"), Blob, 4775, shared.UrlMustParse("https://api.github.com/repos/CDLUC3/mrt-store/git/blobs/d3ae87b904091324ea42af37581d92ee67415143"))

	ee1 := e1.(*entry)
	ee2 := e2.(*entry)

	c.Assert(ee1 == ee2, Equals, true)
}