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
	jobs := node.Jobs
	c.Assert(len(jobs), Equals, 24)

	var job Job
	for _, j := range jobs {
		apiUrl := j.ApiUrl()
		c.Assert(apiUrl.Host, Equals, s.serverUrl.Host)
		if apiUrl.Path == "/job/mrt-store-pub/" {
			job = j
		}
	}
	c.Assert(job, NotNil)
}
