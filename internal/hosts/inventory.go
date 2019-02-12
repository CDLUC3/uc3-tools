package hosts

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

// ------------------------------------------------------------
// Inventory

type Inventory struct {
	invPath string
	hosts []HostRecord
}

func NewInventory(invPath string) *Inventory {
	return &Inventory{invPath: invPath}
}

// ------------------------------
// Exported functions

func (inv *Inventory) Print() error {
	return nil
}

func (inv *Inventory) Hosts() ([]HostRecord, error) {
	if inv.hosts == nil {
		hosts, err := parseInvFile(inv.invPath)
		if err != nil {
			return nil, err
		}
		inv.hosts = hosts
	}
	return inv.hosts, nil
}

// ------------------------------
// Unexported functions

func (inv *Inventory) parse() error {
	return nil
}

// TODO: extract a parsing struct
func parseInvFile(path string) ([]HostRecord, error) {
	invFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := invFile.Close()
		if err != nil {
			err = cerr
		}
	}()

	var recs []HostRecord
	var current *HostRecord
	scanner := bufio.NewScanner(invFile)
	for scanner.Scan() {
		line := scanner.Text()
		isBlank := "" == strings.TrimSpace(line)
		isRecordStart := IsRecordStartLine(line)
		if isBlank || isRecordStart {
			if current != nil {
				sort.Strings(current.CNAMEs)
				recs = append(recs, *current)
				current = nil
			}
		}
		if isBlank {
			continue
		}
		if isRecordStart {
			current, err = NewHostRecord(line)
			if err != nil {
				return nil, err
			}
			continue
		}
		isFqdn := IsFQDNLine(line)
		if isFqdn {
			err = current.AddFQDN(line)
			if err != nil {
				return nil, err
			}
			continue
		}
		_, err = current.MaybeAddCname(line)
		if err != nil {
			return nil, err
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return recs, nil
}

