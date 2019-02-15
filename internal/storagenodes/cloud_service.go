package storagenodes

import (
	"fmt"
	props "github.com/magiconair/properties"
	"net/url"
	"strings"
)

type CloudService struct {
	Name        string
	ServiceType ServiceType
	AccessMode  string
	Endpoint    string
	Credentials CloudCredentials
}

func LoadCloudService(name, svcPropsPath string) (*CloudService, error) {
	svcProps, err := props.LoadFile(svcPropsPath, props.ISO_8859_1)
	if err != nil {
		return nil, err
	}

	serviceType, err := LoadServiceType(svcProps)
	if err != nil {
		return nil, err
	}

	endpoint, err := LoadEndpoint(svcProps)
	if err != nil {
		return nil, err
	}

	credentials, err := LoadCredentials(svcProps)
	if err != nil {
		return nil, err
	}

	node := CloudService{
		Name:        name,
		ServiceType: *serviceType,
		AccessMode:  svcProps.GetString("accessMode", "on-line"),
		Endpoint:    endpoint.String(),
		Credentials: *credentials,
	}
	return &node, nil
}

func LoadEndpoint(svcProps *props.Properties) (*url.URL, error) {
	endpointStr := svcProps.GetString("host", "")
	if endpointStr == "" {
		endpointStr = svcProps.GetString("endPoint", "")
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
