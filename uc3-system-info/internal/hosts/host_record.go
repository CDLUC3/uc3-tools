package hosts

import (
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
	"regexp"
	"sort"
	"strings"
)

const recordStartPrefix = "uc3,"
const fqdnPrefix = "FQDN: "
const cnamePattern = "^\\s +([a-z0-9.-]+cdlib.org)\\.\\s +CNAME.*Private"

var cnameRe = regexp.MustCompile(cnamePattern)

type HostRecord struct {
	Service     string
	Environment string
	Subsystem   string
	Name        string
	FQDN        string
	CNAMEs      []string
}

func (hr *HostRecord) Sprint(format output.Format, hideService bool) (string, error) {
	if hr == nil {
		return "", nil
	}
	fields := []interface{}{
		hr.Service,
		hr.Environment,
		hr.Subsystem,
		hr.Name,
		hr.FQDN,
		hr.CNAMEs,
	}
	if hideService {
		fields = fields[1:]
	}
	return format.Sprint(fields...)
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
	sort.Strings(hr.CNAMEs)
	return nil
}

func (hr *HostRecord) newHostRecord(recordStartLine string) (*HostRecord, error) {
	fields := strings.Split(recordStartLine, ",")
	if len(fields) != 4 {
		return nil, fmt.Errorf("can't parse record start %#v; expected 4 fields, got %d", recordStartLine, len(fields))
	}
	env, svc, name := fields[1], fields[2], fields[3]
	rec := HostRecord{
		Environment: env,
		Service:     svc,
		Name:        name,
		Subsystem:   extractSubsystem(name, svc, env),
	}
	return &rec, nil
}

func (hr *HostRecord) isRecordStart(line string) bool {
	return strings.HasPrefix(line, recordStartPrefix)
}

func (hr *HostRecord) isFQDN(line string) bool {
	return strings.HasPrefix(line, fqdnPrefix)
}

var zoneRe = regexp.MustCompile("2[a-c]")
var numericRe = regexp.MustCompile("[0-9]+")

func extractSubsystem(name string, svc string, env string) string {
	sub := strings.Replace(name, "uc3-", "", -1)
	sub = strings.Replace(sub, svc, "", -1)
	sub = strings.Replace(sub, env, "", -1)
	sub = strings.Replace(sub, "-", "", -1)
	sub = zoneRe.ReplaceAllString(sub, "")
	return numericRe.ReplaceAllString(sub, "")
}
