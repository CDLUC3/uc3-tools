package jenkins

import (
	"github.com/dmolesUC3/mrt-build-info/misc"
	. "gopkg.in/check.v1"
	"io/ioutil"
)

type ConfigSuite struct {}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) TestGoals(c *C) {

	data, err := ioutil.ReadFile("testdata/config.xml")
	c.Assert(err, IsNil)

	url := misc.UrlMustParse("http://builds.cdlib.org/job/mrt-store-pub/config.xml")
	config, err := ConfigFromBytes(data, url)
	c.Assert(err, IsNil)

	expected := "clean install -DpropertyDir=$propertyDirName -Djava.compiler=1.8"
	goals, err := config.Goals()
	c.Assert(err, IsNil)
	c.Assert(goals, Equals, expected)
}