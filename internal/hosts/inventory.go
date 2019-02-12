package hosts

// ------------------------------------------------------------
// Inventory

type Inventory struct {
	invPath string
}

func NewInventory(invPath string) *Inventory {
	return &Inventory{invPath: invPath}
}

// ------------------------------
// Exported functions

func (inv *Inventory) Print() error {
	return nil
}

// ------------------------------
// Unexported functions

func (inv *Inventory) parse() error {
	return nil
}