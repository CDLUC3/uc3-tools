package hosts

import (
	"fmt"
	"regexp"
	"strings"
)

const recordStartPrefix = "uc3,"
const fqdnPrefix = "FQDN: "
const cnamePattern = "^\\s +([a-z0-9.-]+cdlib.org)\\.\\s +CNAME.*Private"

var cnameRe = regexp.MustCompile(cnamePattern)

type HostRecord struct {
	Service     string
	Name        string
	FQDN        string
	Environment string
	CNAMEs      []string
}

func NewHostRecord(recordStartLine string) (*HostRecord, error) {
	if !IsRecordStartLine(recordStartLine) {
		return nil, fmt.Errorf("not a record start line: %#v", recordStartLine)
	}
	fields := strings.Split(recordStartLine, ",")
	if len(fields) != 4 {
		return nil, fmt.Errorf("can't parse record start %#v; expected 4 fields, got %d", recordStartLine, len(fields))
	}
	env, svc, name := fields[0], fields[1], fields[2]
	rec := HostRecord{
		Environment: env,
		Service:     svc,
		Name:        name,
	}
	return &rec, nil
}

func (hr *HostRecord) ParseLine(line string) (*HostRecord, error) {
	if "" == strings.TrimSpace(line) {
		return hr, nil
	}
	if IsRecordStartLine(line) {
		return NewHostRecord(line)
	}
	if hr == nil {
		return nil, fmt.Errorf("can't parse %#v into nil record", line)
	}
	if IsFQDNLine(line) {
		err := hr.AddFQDN(line)
		return hr, err
	}
	err := hr.MaybeAddCname(line)
	return hr, err
}

func (hr *HostRecord) AddFQDN(fqdnLine string) error {
	if hr == nil {
		return fmt.Errorf("can't add FQDN from %#v to nil record", fqdnLine)
	}
	if !IsFQDNLine(fqdnLine) {
		return fmt.Errorf("not a FQDN line: %#v", fqdnLine)
	}
	hr.FQDN = strings.TrimPrefix(fqdnLine, fqdnPrefix)
	return nil
}

func (hr *HostRecord) MaybeAddCname(line string) (error) {
	cnameMatch := cnameRe.FindStringSubmatch(line)
	if cnameMatch == nil {
		return nil
	}
	if hr == nil {
		return fmt.Errorf("can't add CNAME from %#v to nil record", line)
	}
	hr.CNAMEs = append(hr.CNAMEs, cnameMatch[1])
	return nil
}

func IsRecordStartLine(line string) bool {
	return strings.HasPrefix(line, recordStartPrefix)
}

func IsFQDNLine(line string) bool {
	return strings.HasPrefix(line, fqdnPrefix)
}

