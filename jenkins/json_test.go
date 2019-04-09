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

func (s *JsonSuite) TestParseJob(c *C) {
	data, _ := ioutil.ReadFile("testdata/job.json")
	var job Job
	err := json.Unmarshal(data, &job)
	c.Assert(err, IsNil)

	build := job.LastSuccessfulBuild
	c.Assert(build.Number, Equals, 93)

	c.Assert(build.URL, Equals, "http://builds.cdlib.org/job/mrt-store-pub/93/")
}

func (s *JsonSuite) TestParseBuild(c *C) {
	data, _ := ioutil.ReadFile("testdata/build.json")
	var build Build
	err := json.Unmarshal(data, &build)
	c.Assert(err, IsNil)

	actions := build.Actions
	c.Assert(len(actions), Equals, 9)

	bdAction := actions[2]
	c.Assert(bdAction.Class, Equals, "hudson.plugins.git.util.BuildData")

	rURLs := bdAction.RemoteURLs
	c.Assert(len(rURLs), Equals, 1)
	c.Assert(rURLs[0], Equals, "https://github.com/CDLUC3/mrt-store.git")

	rev := bdAction.LastBuiltRevision
	c.Assert(rev, NotNil)
	c.Assert(rev.SHA1, Equals, "af174ac555758a1c639a7a3da39e022d9fdbf3a6")

	branches := rev.Branches
	c.Assert(len(branches), Equals, 1)

	branch := branches[0]
	c.Assert(branch, NotNil)
	c.Assert(branch.SHA1, Equals, "af174ac555758a1c639a7a3da39e022d9fdbf3a6")
	c.Assert(branch.Name, Equals, "refs/remotes/origin/master")
}