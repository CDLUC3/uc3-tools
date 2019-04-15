package jenkins

import (
	"encoding/json"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"regexp"
)

var _ = Suite(&JsonSuite{})

type JsonSuite struct{}

func (s *JsonSuite) TestParseNode(c *C) {
	data, _ := ioutil.ReadFile("testdata/node.json")
	var node Node = &node{}
	err := json.Unmarshal(data, node)
	c.Assert(err, IsNil)

	jobs := node.Jobs()
	c.Assert(len(jobs), Equals, 24)

	expectedNames := []string{
		"cdl-zk-queue",
		"git-core2",
		"git-dataONE",
		"git-mrt-xoai",
		"git-mrt-zoo",
		"Merritt Development Submission (Full Stack Test)",
		"Merritt Production Submission (Full Stack Test)",
		"Merritt Stage Submission (Full Stack Test)",
		"mrt-build-audit",
		"mrt-build-inv",
		"mrt-build-mysql",
		"mrt-build-oai",
		"mrt-build-replic",
		"mrt-build-s3",
		"mrt-build-store",
		"mrt-build-sword",
		"mrt-cloudhost-pub",
		"mrt-ingest-dev",
		"mrt-ingest-stage",
		"mrt-jetty-cloudhost",
		"mrt-s3-pub",
		"mrt-store-pub",
		"mrt-test",
		"test-gittest",
	}
	for i, j := range jobs {
		c.Check(j.Name(), Equals, expectedNames[i])
	}
}

func (s *JsonSuite) TestParseJob(c *C) {
	data, _ := ioutil.ReadFile("testdata/job.json")
	var job Job = &job{}
	err := json.Unmarshal(data, job)
	c.Assert(err, IsNil)

	c.Assert(job.Name(), Equals, "mrt-store-pub")

	build, err := job.LastSuccess()
	c.Assert(err, IsNil)

	c.Assert(build.BuildNumber(), Equals, 93)
}

func (s *JsonSuite) TestParseBuild(c *C) {
	data, _ := ioutil.ReadFile("testdata/build.json")
	var build Build = &build{}
	err := json.Unmarshal(data, build)
	c.Assert(err, IsNil)

	c.Assert(build.BuildNumber(), Equals, 93)
	scmUrl, err := build.SCMUrl()
	c.Assert(err, IsNil)
	c.Assert(scmUrl, Equals, "https://github.com/CDLUC3/mrt-store.git")

	sha1, err := build.SHA1()
	c.Assert(err, IsNil)
	c.Assert(sha1, Equals, "af174ac555758a1c639a7a3da39e022d9fdbf3a6")

	artifacts, err := build.Artifacts()
	c.Assert(err, IsNil)

	c.Assert(len(artifacts), Equals, 3)

	expectedArtifacts := []string{"mrt-storepub", "mrt-storepub-src", "mrt-storewar"}
	expectedTypes := []string{"pom", "jar", "war"}

	for i, a := range artifacts {
		c.Check(a.GroupId(), Equals, "org.cdlib.mrt")
		c.Check(a.Version(), Equals, "1.0-SNAPSHOT")

		artifact := expectedArtifacts[i]
		_type := expectedTypes[i]
		c.Check(a.ArtifactId(), Equals, artifact)
		c.Check(a.Packaging(), Equals, _type)
	}
}

func (s JsonSuite) TestBuildCommit(c *C) {
	re := regexp.MustCompile("([^/]+)/([^@]+)@([a-f0-9]+)")
	repoByFile := map[string]string {
		"testdata/build.json": "CDLUC3/mrt-store@af174ac555758a1c639a7a3da39e022d9fdbf3a6",
		"testdata/build-private.json": "cdlib/mrt-conf-prv@691770b9182d3870c85aba8ca776c0a3e85aa57e",
	}
	for file, expected := range repoByFile {
		data, _ := ioutil.ReadFile(file)
		var build Build = &build{}
		err := json.Unmarshal(data, build)
		c.Assert(err, IsNil)

		owner, repo, sha1, err := build.Commit()
		c.Assert(err, IsNil)

		matches := re.FindStringSubmatch(expected)
		c.Assert(owner, Equals, matches[1])
		c.Assert(repo, Equals, matches[2])
		c.Assert(sha1, Equals, matches[3])
	}
}

