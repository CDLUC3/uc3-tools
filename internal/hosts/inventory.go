package hosts

import (
	"bufio"
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/outputfmt"
	"os"
	"sort"
)

// ------------------------------------------------------------
// Inventory

type Inventory struct {
	Hosts []HostRecord
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

func (inv *Inventory) Print(format outputfmt.Format) {
	for _, h := range inv.Hosts {
		fmt.Printf("%v%v%v\n",
			format.Prefix(),
			// TODO: align columns in markdown format
			h.ToDelimitedString(
				format.FieldSeparator(), format.InnerSeparator(),
			),
			format.Suffix(),
		)
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
