package jenkins

import (
	"encoding/json"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/git"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"regexp"
	"sort"
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

	configUrl := job.ConfigUrl()
	c.Assert(configUrl, NotNil)

	expectedConfigUrl := "http://builds.cdlib.org/job/mrt-store-pub/config.xml"
	c.Assert(configUrl.String(), Equals, expectedConfigUrl)

	build, err := job.LastSuccess()
	c.Assert(err, IsNil)

	c.Assert(build.BuildNumber(), Equals, 93)

	params := job.Parameters()
	c.Assert(len(params), Equals, 1)

	param := params[0]
	c.Check(param.Name(), Equals, "propertyDirName")
	c.Check(param.Default(), Equals, "dev")

	expected := []string{"dev", "stage", "prod", "audit", "test"}
	sort.Strings(expected)
	choices := param.Choices()
	c.Assert(len(choices), Equals, len(expected))
	for i, actual := range choices {
		c.Check(actual, Equals, expected[i])
	}
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
	c.Assert(sha1, Equals, git.SHA1("af174ac555758a1c639a7a3da39e022d9fdbf3a6"))

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
	repoByFile := map[string]string{
		"testdata/build.json":         "CDLUC3/mrt-store@af174ac555758a1c639a7a3da39e022d9fdbf3a6",
		"testdata/build-private.json": "cdlib/mrt-conf-prv@691770b9182d3870c85aba8ca776c0a3e85aa57e",
		"testdata/build-ingest-dev.json": "CDLUC3/mrt-ingest@a572b8c116ef56936ffd866749bd667d29a36fdd",
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
		c.Assert(sha1, Equals, git.SHA1(matches[3]))
	}
}
