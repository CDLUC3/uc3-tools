package hosts

import (
	"bufio"
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
	"os"
	"sort"
	"time"
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
	sort.Sort(ByServiceAndEnvironment(hosts))
	return &Inventory{Hosts: hosts}, nil
}

// ------------------------------
// Exported functions

func (inv *Inventory) Print(format output.Format, header bool, footer bool, service string) error {
	hideService := service != ""

	if header {
		headerFields := []string{
			"Service",
			"Environment",
			"Subsystem",
			"Name",
			"FQDN",
			"CNAMEs",
		}
		if hideService {
			headerFields = headerFields[1:]
		}
		fmt.Print(format.SprintHeader(headerFields...))
	}

	// TODO: "pretty" TSV:
	//   - 1 line per CNAME
	//   - don't repeat left columns

	for _, h := range inv.Hosts {
		if hideService && h.Service != service {
			continue
		}
		line, err := h.Sprint(format, hideService)
		if err != nil {
			return err
		}
		fmt.Println(line)
	}

	if footer {
		fmt.Printf("\nGenerated %v\n", time.Now().Format(time.RFC3339))
	}

	return nil
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
