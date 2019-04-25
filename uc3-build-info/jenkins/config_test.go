package jenkins

import (
	"github.com/CDLUC3/uc3-tools/mrt-build-info/shared"
	. "gopkg.in/check.v1"
	"io/ioutil"
)

type ConfigSuite struct{}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) TestGoals(c *C) {
	data, err := ioutil.ReadFile("testdata/config.xml")
	c.Assert(err, IsNil)

	url := shared.UrlMustParse("http://builds.cdlib.org/job/mrt-store-pub/config.xml")
	config, err := ConfigFromBytes(data, url)
	c.Assert(err, IsNil)

	expected := "clean install -DpropertyDir=$propertyDirName -Djava.compiler=1.8"
	goals := config.Goals()
	c.Assert(goals, Equals, expected)
}

func (s *ConfigSuite) TestMavenParamToValue(c *C) {
	data, err := ioutil.ReadFile("testdata/config.xml")
	c.Assert(err, IsNil)

	url := shared.UrlMustParse("http://builds.cdlib.org/job/mrt-store-pub/config.xml")
	config, err := ConfigFromBytes(data, url)
	c.Assert(err, IsNil)

	mavenParameters := config.MavenParamToValue()
	c.Assert(mavenParameters["propertyDir"], Equals, "$propertyDirName")
	c.Assert(mavenParameters["java.compiler"], Equals, "1.8")
}

func (s *ConfigSuite) TestBuildRootEmpty(c *C) {
	data, err := ioutil.ReadFile("testdata/config.xml")
	c.Assert(err, IsNil)

	url := shared.UrlMustParse("http://builds.cdlib.org/job/mrt-store-pub/config.xml")
	config, err := ConfigFromBytes(data, url)
	c.Assert(err, IsNil)

	c.Assert(config.BuildRoot(), Equals, "")
}

func (s *ConfigSuite) TestBuildRoot(c *C) {
	data, err := ioutil.ReadFile("testdata/config-xml/mrt-build-store-config.xml")
	c.Assert(err, IsNil)

	url := shared.UrlMustParse("http://builds.cdlib.org/job/mrt-build-store/config.xml")
	config, err := ConfigFromBytes(data, url)
	c.Assert(err, IsNil)

	c.Assert(config.BuildRoot(), Equals, "store-conf")
}
