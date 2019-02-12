package hosts

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// ------------------------------------------------------------
// Inventory

type Inventory struct {
	Hosts   []HostRecord
}

func NewInventory(invPath string) (*Inventory, error) {
	hosts, err := parseInvFile(invPath)
	if err != nil {
		return nil, err
	}
	sort.Sort(BySvcEnvNameAndFQDN(hosts))
	return &Inventory{Hosts: hosts}, nil
}

// ------------------------------
// Exported functions

func (inv *Inventory) Print() {
	for _, h := range inv.Hosts {
		fmt.Println(h.ToDelimitedString("\t", ", "))
	}
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
