package storage

import (
	"fmt"
	props "github.com/magiconair/properties"
)

type ServiceType int

const (
	unknown ServiceType = iota
	s3
	swift
	cloudhost
	pairtree
)

func (s ServiceType) String() string {
	switch (s) {
	case s3:
		return "s3"
	case swift:
		return "swift"
	case cloudhost:
		return "cloudhost"
	case pairtree:
		return "pairtree"
	default:
		return "unknown"
	}
}

func LoadServiceType(nodeProps *props.Properties) (*ServiceType, error) {
	serviceTypeStr := nodeProps.GetString("serviceType", "")
	serviceType, err := ParseServiceType(serviceTypeStr)
	if err != nil {
		return nil, err
	}
	return &serviceType, nil
}

func ParseServiceType(serviceTypeStr string) (ServiceType, error) {
	if "aws" == serviceTypeStr || "minio" == serviceTypeStr {
		return s3, nil
	}
	if "swift" == serviceTypeStr {
		return swift, nil
	}
	if "cloudhost" == serviceTypeStr {
		return cloudhost, nil
	}
	if "pairtree" == serviceTypeStr {
		return pairtree, nil
	}
	return unknown, fmt.Errorf("can't parse service type %#v", serviceTypeStr)
}
