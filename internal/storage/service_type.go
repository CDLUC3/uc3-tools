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

func (s ServiceType) CLIExampleFile(credentials *CloudCredentials, endpoint, container, key string) (string, error) {
	if s == s3 {
		return s3Example(credentials, endpoint, container, key), nil
	}
	if s == swift {
		return swiftExampleFile(credentials, endpoint, container, key), nil
	}
	return s.cliExampleNotAvailable()

}

func (s ServiceType) CLIExampleObject(credentials *CloudCredentials, endpoint, container, ark string) (string, error) {
	if s == s3 {
		return s3Example(credentials, endpoint, container, ark), nil
	}
	if s == swift {
		return swiftExampleObject(credentials, endpoint, container, ark), nil
	}
	return s.cliExampleNotAvailable()
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

// ------------------------------------------------------------
// Unexported symbols

func (s ServiceType) cliExampleNotAvailable() (string, error) {
	return "", fmt.Errorf("(command-line example not available for %s service)", s.String())
}

func s3Example(credentials *CloudCredentials, endpoint, container, key string) string {
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
	return strings.Join(cmd, " ")
}

func swiftExampleFile(credentials *CloudCredentials, endpoint, container, key string) string {
	var cmd []string
	cmd = append(cmd, swiftCmd())
	if credentials != nil {
		cmd = append(cmd, "-A", shellSafe(endpoint), "-U", shellSafe(credentials.Key), "-K", shellSafe(credentials.Secret))
	}
	cmd = append(cmd, "stat", shellSafe(container), shellSafe(key))
	return strings.Join(cmd, " ")
}

func swiftExampleObject(credentials *CloudCredentials, endpoint, container, ark string) string {
	var cmd []string
	cmd = append(cmd, swiftCmd())
	if credentials != nil {
		cmd = append(cmd, "-A", shellSafe(endpoint), "-U", shellSafe(credentials.Key), "-K", shellSafe(credentials.Secret))
	}
	cmd = append(cmd, "list", shellSafe(container), "--prefix", shellSafe(ark))
	return strings.Join(cmd, " ")
}

func swiftCmd() string {
	if runtime.GOOS == "darwin" {
		// on macOS, "swift" is likely /usr/bin/swift, the Swift compiler; this is
		// our best guess at the default location of the OpenStack Swift CLI
		return "/usr/local/bin/swift"
	}
	return "swift"
}

func shellSafe(s string) string {
	// single-quote the string for safety
	return fmt.Sprintf("'%s'", strings.Replace(s, "'", "'\\''", -1))
}

