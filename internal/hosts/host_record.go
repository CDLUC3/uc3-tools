package hosts

type HostRecord struct {
	FQDN string
	Environment string
	CNAMEs []string
}