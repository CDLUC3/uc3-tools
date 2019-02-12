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

func (hr *HostRecord) ParseLine(line string) (*HostRecord, error) {
	if "" == strings.TrimSpace(line) {
		return hr, nil
	}
	if hr.isRecordStart(line) {
		return hr.newHostRecord(line)
	}
	if hr == nil {
		return nil, fmt.Errorf("can't parse %#v into nil record", line)
	}
	if hr.isFQDN(line) {
		err := hr.addFQDN(line)
		return hr, err
	}
	err := hr.maybeAddCname(line)
	return hr, err
}

func (hr *HostRecord) addFQDN(fqdnLine string) error {
	if hr == nil {
		return fmt.Errorf("can't add FQDN from %#v to nil record", fqdnLine)
	}
	hr.FQDN = strings.TrimPrefix(fqdnLine, fqdnPrefix)
	return nil
}

func (hr *HostRecord) maybeAddCname(line string) error {
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

func (hr *HostRecord) newHostRecord(recordStartLine string) (*HostRecord, error) {
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

func (hr *HostRecord) isRecordStart(line string) bool {
	return strings.HasPrefix(line, recordStartPrefix)
}

func (hr *HostRecord) isFQDN(line string) bool {
	return strings.HasPrefix(line, fqdnPrefix)
}

