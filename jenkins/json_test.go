package jenkins

import (
	"encoding/json"
	. "gopkg.in/check.v1"
	"io/ioutil"
)

var _ = Suite(&JsonSuite{})

type JsonSuite struct {}

func (s *JsonSuite) TestParseNode(c *C) {
	data, _ := ioutil.ReadFile("testdata/node.json")
	var node Node
	err := json.Unmarshal(data, &node)
	c.Assert(err, IsNil)

	jobs := node.Jobs
	c.Assert(len(jobs), Equals, 24)
}