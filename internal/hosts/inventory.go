package hosts

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

func parseInvFile(path string) ([]HostRecord, error) {
	return nil, nil
}