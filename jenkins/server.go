package jenkins

import (
	"fmt"
	"net/url"
)

// JenkinsServer represents a Jenkins server
type JenkinsServer interface {
	Node() (Node, error)
}

func DefaultServer() JenkinsServer {
	server, err := ServerFromUrl("http://builds.cdlib.org/")
	if err == nil {
		return server
	}
	// should never happen
	panic(err)
}

func ServerFromUrl(urlStr string) (JenkinsServer, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, fmt.Errorf("server URL '%v' is not absolute", urlStr)
	}
	if !isApiUrl(u) {
		u = toApiUrl(u)
	}
	return &jenkinsServer{apiRoot: u}, nil
}

// ------------------------------------------------------------
// Unexported symbols

type jenkinsServer struct {
	apiRoot *url.URL
	node    Node
}

func (s *jenkinsServer) Node() (Node, error) {
	if s.node == nil {
		var n Node = &node{}
		if err := unmarshal(s.apiRoot, n); err != nil {
			return nil, err
		}
		s.node = n
	}
	return s.node, nil
}
