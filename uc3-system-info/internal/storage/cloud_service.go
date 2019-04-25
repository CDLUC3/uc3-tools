package storage

import (
	"fmt"
	"github.com/CDLUC3/uc3-tools/uc3-system-info/internal/output"
	props "github.com/magiconair/properties"
	"net/url"
	"strings"
)

type CloudService struct {
	Name        string
	ServiceType ServiceType
	AccessMode  string
	Endpoint    string
	Credentials *CloudCredentials
}

func (c *CloudService) Sprint(format output.Format) string {
	if c == nil {
		return ""
	}
	str, _ := format.Sprint(c.Name, c.ServiceType, c.AccessMode, c.Endpoint, c.Key(), c.Secret())
	return str
}

func (c *CloudService) Key() string {
	creds := c.Credentials
	if creds == nil {
		return ""
	}
	return creds.Key
}

func (c *CloudService) Secret() string {
	creds := c.Credentials
	if creds == nil {
		return ""
	}
	return creds.Secret
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

	endpoint, err := LoadEndpoint(svcProps, serviceType)
	if err != nil {
		return nil, err
	}

	credentials, err := LoadCredentials(svcProps)
	if err != nil {
		return nil, err
	}

	var endpointStr string
	if endpoint == nil {
		endpointStr = ""
	} else {
		endpointStr = endpoint.String()
	}

	node := CloudService{
		Name:        name,
		ServiceType: *serviceType,
		AccessMode:  svcProps.GetString("accessMode", "on-line"),
		Endpoint:    endpointStr,
		Credentials: credentials,
	}
	return &node, nil
}

func LoadEndpoint(svcProps *props.Properties, serviceType *ServiceType) (*url.URL, error) {
	endpointStr := svcProps.GetString("host", "")
	if endpointStr == "" {
		endpointStr = svcProps.GetString("endPoint", "")
	}
	// TODO: move this to ServiceType
	if serviceType != nil && *serviceType == swift {
		endpointStr = fmt.Sprintf("%s/auth/v1.0", endpointStr)
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
