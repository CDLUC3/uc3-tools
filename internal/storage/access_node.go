package storage

import (
	"fmt"
	props "github.com/magiconair/properties"
	"net/url"
	"strings"
)

type AccessNode struct {
	ServiceType ServiceType
	AccessMode  string
	Endpoint    *url.URL
	Credentials *CloudCredentials
}

func LoadAccessNode(nodePropsPath string) (*AccessNode, error) {
	nodeProps, err := props.LoadFile(nodePropsPath, props.ISO_8859_1)
	if err != nil {
		return nil, err
	}

	serviceType, err := LoadServiceType(nodeProps)
	if err != nil {
		return nil, err
	}

	endpoint, err := LoadEndpoint(nodeProps)
	if err != nil {
		return nil, err
	}

	credentials, err := LoadCredentials(nodeProps)
	if err != nil {
		return nil, err
	}

	node := AccessNode{
		ServiceType: *serviceType,
		AccessMode:  nodeProps.GetString("accessMode", "on-line"),
		Endpoint: endpoint,
		Credentials: credentials,
	}
	return &node, nil
}

func LoadEndpoint(nodeProps *props.Properties) (*url.URL, error) {
	endpointStr := nodeProps.GetString("host", "")
	if endpointStr == "" {
		endpointStr = nodeProps.GetString("endPoint", "")
	}
	return toAbsoluteUrl(endpointStr)
}

func toAbsoluteUrl(endpointStr string) (*url.URL, error) {
	if endpointStr == "" {
		return nil, nil
	}
	endpointUrl, err := url.Parse(endpointStr)
	if err != nil {
		return nil, err
	}
	if endpointUrl.IsAbs() {
		return endpointUrl, nil
	}

	pathComponents := strings.Split(endpointUrl.Path, "/")
	host := pathComponents[0]

	absUrlStr := fmt.Sprintf("https://%v:443", host)
	if len(pathComponents) > 1 {
		path := pathComponents[1:]
		absUrlStr = absUrlStr + "/" + strings.Join(path, "/")
	}
	return url.Parse(absUrlStr)
}