package jenkins

import (
	"fmt"
	"net/http"
	"net/url"
)

// JenkinsServer represents a Jenkins server
type JenkinsServer interface {
	Node() (*Node, error)
}

func DefaultServer() JenkinsServer {
	server, err := ServerFromUrl("http://builds.cdlib.org/")
	if err == nil {
		return server
	}
	panic(err)
}

func ServerFromUrl(urlStr string) (JenkinsServer, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, fmt.Errorf("jenkinsServer URL '%v' is not absolute", urlStr)
	}
	return &jenkinsServer{serverUrl: u}, nil
}

// ------------------------------------------------------------
// Unexported symbols

type jenkinsServer struct {
	serverUrl *url.URL
	client *http.Client
	node *Node
}

func (s *jenkinsServer) Node() (*Node, error) {
	if s.node == nil {
		s.node = &Node{}
	}
	err := unmarshal(s.apiRoot(), s.node)
	if err != nil {
		return nil, err
	}

	for _, job := range s.node.Jobs {
		err = job.load()
		if err != nil {
			break
		}
	}

	return s.node, err
}

func (s *jenkinsServer) apiRoot() *url.URL {
	return toApiUrl(s.serverUrl)
}

