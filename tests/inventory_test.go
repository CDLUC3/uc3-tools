package tests

import (
	. "github.com/dmolesUC3/mrt-system-info/internal/hosts"
	. "gopkg.in/check.v1"
)

var _ = Suite(&InventorySuite{})

type InventorySuite struct {
}

func (s *InventorySuite) TestHosts(c *C) {
	inv := NewInventory("testdata/uc3-inventory.txt")
	hosts, err := inv.Hosts()
	c.Assert(err, IsNil)
	c.Assert(hosts, NotNil)
	c.Assert(len(hosts), Equals, 80)
}