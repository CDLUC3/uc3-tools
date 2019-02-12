package hosts

import (
	"bufio"
	"os"
	"sort"
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
		sort.Sort(BySvcEnvNameAndFQDN(hosts))
		inv.hosts = hosts
	}
	return inv.hosts, nil
}

// ------------------------------
// Unexported functions

func (inv *Inventory) parse() error {
	return nil
}

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
		hr, err := current.ParseLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		if hr != current {
			if current != nil {
				recs = append(recs, *current)
			}
			current = hr
		}
	}
	// TODO: find a way not to repeat this (use last slice element?)
	if current != nil {
		recs = append(recs, *current)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return recs, nil
}

