package tests

import (
	. "github.com/dmolesUC3/mrt-system-info/internal/hosts"
	. "gopkg.in/check.v1"
)

var _ = Suite(&InventorySuite{})

type InventorySuite struct {
}

var noCnameHosts = []string {
	"uc3-mrtdat1-dev.cdlib.org",
	"uc3-mrtdocker2-stg.cdlib.org",
	"uc3-ldap-prd-2a.cdlib.org",
	"uc3-dryadsolr-stg.cdlib.org",
	"uc3-dryadsolr-dev.cdlib.org",
	"uc3-ldapvm-stg.cdlib.org",
	"uc3-mrtsandbox2-stg.cdlib.org",
	"uc3-dryadui-stg-2c.cdlib.org",
	"uc3-ldap-prd-2c.cdlib.org",
}

func isNoCnameHost(fqdn string) bool {
	for _, h := range noCnameHosts {
		if h == fqdn {
			return true
		}
	}
	return false
}

func (s *InventorySuite) TestHosts(c *C) {
	inv, err := NewInventory("testdata/uc3-inventory.txt")
	c.Assert(err, IsNil)
	hosts := inv.Hosts
	c.Assert(hosts, NotNil)
	c.Assert(len(hosts), Equals, 80)

	for _, h := range hosts {
		c.Assert(h.Service, NotNil)
		c.Assert(h.Name, NotNil)
		c.Assert(h.FQDN, NotNil)
		c.Assert(h.Environment, NotNil)
		if isNoCnameHost(h.FQDN) {
			c.Assert(len(h.CNAMEs), Equals, 0, Commentf("unexpected CNAMEs found for %v: %#v", h.FQDN, h.CNAMEs))
		} else {
			c.Assert(len(h.CNAMEs) > 0, Equals, true, Commentf("no CNAMEs found for %v", h.FQDN))
		}
	}
}