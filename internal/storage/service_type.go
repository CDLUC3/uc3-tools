package storage

import (
	"fmt"
	props "github.com/magiconair/properties"
)

type ServiceType int

const (
	s3 ServiceType = iota
	swift
	unknown
)

func LoadServiceType(nodeProps *props.Properties) (*ServiceType, error) {
	serviceTypeStr := nodeProps.GetString("serviceTypeStr", "")
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
	return unknown, fmt.Errorf("can't parse service type %#v", serviceTypeStr)
}
