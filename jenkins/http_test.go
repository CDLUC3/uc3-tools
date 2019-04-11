package jenkins

import (
	"fmt"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

var _ = Suite(&HttpSuite{})

type HttpSuite struct {
	server    *httptest.Server
	serverUrl *url.URL
	jenkins   JenkinsServer
}

// ------------------------------------------------------------
// Fixture

var data = map[string]string{
	"/api/json":                      "testdata/node.json",
	"/job/mrt-store-pub/api/json":    "testdata/job.json",
	"/job/mrt-store-pub/93/api/json": "testdata/build.json",
}

func (s *HttpSuite) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if file, ok := data[r.URL.Path]; ok {
		body, err := ioutil.ReadFile(file)
		if err == nil {
			// note that URL.Host includes port
			body = []byte(strings.ReplaceAll(string(body), "builds.cdlib.org", s.serverUrl.Host))
			n, err := w.Write(body)
			if err != nil {
				panic(err)
			}
			if n != len(body) {
				panic(fmt.Errorf("wrong number of bytes: wrote %d, expected %d", n, len(body)))
			}
			return
		}
	}
	w.WriteHeader(500)
}

func (s *HttpSuite) SetUpTest(c *C) {
	s.server = httptest.NewServer(http.HandlerFunc(s.HandleRequest))
	s.serverUrl = urlMustParse(s.server.URL)
	jenkins, err := ServerFromUrl(s.server.URL)
	if err != nil {
		c.Error(err)
		c.FailNow()
	}
	s.jenkins = jenkins
}

// ------------------------------------------------------------
// Tests

func (s *HttpSuite) TestNode(c *C) {
	node, err := s.jenkins.Node()
	c.Assert(err, IsNil)
	c.Assert(node, NotNil)
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

	var job Job
	for _, j := range jobs {
		if j.Name() == "mrt-store-pub" {
			job = j
			break
		}
	}
	c.Assert(job, NotNil)

	build, err := job.LastSuccess()
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
		c.Check(a.Group(), Equals, "org.cdlib.mrt")
		c.Check(a.Version(), Equals, "1.0-SNAPSHOT")

		artifact := expectedArtifacts[i]
		_type := expectedTypes[i]
		c.Check(a.Artifact(), Equals, artifact)
		c.Check(a.Type(), Equals, _type)

		expectedFile := fmt.Sprintf("%v-1.0-SNAPSHOT.%v", artifact, _type)
		c.Check(a.File(), Equals, expectedFile)
	}
}
