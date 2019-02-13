package hosts

import "strings"

type ByServiceAndEnvironment []HostRecord

func (r ByServiceAndEnvironment) Len() int {
	return len(r)
}

func (r ByServiceAndEnvironment) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ByServiceAndEnvironment) Less(i, j int) bool {
	return r.compare(r[i], r[j]) < 0
}

func (r ByServiceAndEnvironment) compare(r1, r2 HostRecord) int {
	for _, c := range []func() int {
		func() int { return strings.Compare(r1.Service, r2.Service) },
		func() int { return compareEnvs(r1.Environment, r2.Environment) },
		func() int { return strings.Compare(r1.Subsystem, r2.Subsystem) },
		func() int { return strings.Compare(r1.Name, r2.Name)},
		func() int { return strings.Compare(r1.FQDN, r2.FQDN)},
	} {
		if order := c(); order != 0 {
			return order
		}
	}
	return 0
}

func compareEnvs(e1, e2 string) int {
	if e1 == "dev" {
		if e2 == "dev" {
			return 0
		}
		return -1
	}
	if e1 == "stg" {
		if e2 == "stg" {
			return 0
		}
		if e2 == "dev" {
			return 1
		}
		return -1
	}
	if e1 == "prd" {
		if e2 == "prd" {
			return 0
		}
		return 1
	}
	return strings.Compare(e1, e2)
}

