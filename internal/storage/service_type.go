package storage

import (
	"crypto/md5"
	"fmt"
	props "github.com/magiconair/properties"
	"runtime"
	"strings"
)

type ServiceType int

const (
	unknown ServiceType = iota
	s3
	swift
	cloudhost
	pairtree

	arkMd5Suffix = "<ark-md5-prefix>"
)

func (s ServiceType) String() string {
	switch s {
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

func (s ServiceType) ContainerFor(containerBase, ark string) string {
	if s == swift && strings.HasSuffix(containerBase, "__") {
		hash := md5.New()
		hash.Write(([]byte)(ark))
		resultStr := fmt.Sprintf("%x", hash.Sum(nil))
		return containerBase + resultStr[0:3]
	}
	return containerBase
}

func (s ServiceType) ContainerGeneric(containerBase string) string {
	if s == swift && strings.HasSuffix(containerBase, "__") {
		return containerBase + arkMd5Suffix
	}
	return containerBase
}

func (s ServiceType) CLIExample(credentials *CloudCredentials, endpoint, container, key string) (string, error) {
	shellSafe := func(s string) string {
		// single-quote the string for safety
		return fmt.Sprintf("'%s'", strings.Replace(s, "'", "'\\''", -1))
	}

	if s == s3 {
		var cmd []string
		if credentials != nil {
			cmd = append(cmd, fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", shellSafe(credentials.Key)))
			cmd = append(cmd, fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", shellSafe(credentials.Secret)))
		}
		cmd = append(cmd, "aws")
		if endpoint != "" {
			cmd = append(cmd, "--endpoint", shellSafe(endpoint))
		}
		cmd = append(cmd, "s3", "ls")
		s3Url := fmt.Sprintf("s3://%s/%s", container, key)
		cmd = append(cmd, shellSafe(s3Url))
		return strings.Join(cmd, " "), nil
	}
	if s == swift {
		var cmd []string
		swiftCmd := "swift"
		if runtime.GOOS == "darwin" {
			// on macOS, "swift" is likely /usr/bin/swift, the Swift compiler; this is
			// our best guess at the default location of the OpenStack Swift CLI
			swiftCmd = "/usr/local/bin/swift"
		}
		cmd = append(cmd, swiftCmd)
		if credentials != nil {
			cmd = append(cmd, "-A", shellSafe(endpoint), "-U", shellSafe(credentials.Key), "-K", shellSafe(credentials.Secret))
		}
		cmd = append(cmd, "stat", shellSafe(container), shellSafe(key))
		return strings.Join(cmd, " "), nil
	}
	return "", fmt.Errorf("(command-line example not available for %s service)", s.String())
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
