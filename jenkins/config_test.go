package jenkins

import (
	"github.com/beevik/etree"
	"github.com/dmolesUC3/mrt-build-info/misc"
	. "gopkg.in/check.v1"
)

type ConfigSuite struct {}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) TestGoals(c *C) {
	doc := etree.NewDocument()
	// TODO: figure out how to parse XML 1.1 :P
	err := doc.ReadFromFile("testdata/config.xml")
	c.Assert(err, IsNil)

	var config Config = &config{
		url: misc.UrlMustParse("http://builds.cdlib.org/job/mrt-store-pub/config.xml"),
		doc: doc,
	}

	expected := "clean install -DpropertyDir=$propertyDirName -Djava.compiler=1.8"
	goals, err := config.Goals()
	c.Assert(err, IsNil)
	c.Assert(goals, Equals, expected)
}