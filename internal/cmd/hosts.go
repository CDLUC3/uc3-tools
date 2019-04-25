package cmd

import (
	. "github.com/dmolesUC3/uc3-system-info/internal/hosts"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
	"github.com/spf13/cobra"
)

type HostFlags struct {
	FormatStr string
	Header    bool
	Footer    bool
	service string
}

func (f *HostFlags) PrintInventory(invPath string) error {
	format, err := output.ToFormat(f.FormatStr)
	if err != nil {
		return err
	}
	inv, err := NewInventory(invPath)
	if err != nil {
		return err
	}
	return inv.Print(format, f.Header, f.Footer, f.service)
}

const (
	hostsExamples = `
		uc3-system-info hosts uc3-inventory.txt
		uc3-system-info hosts uc3-inventory.txt --format md --header --footer
		uc3-system-info hosts uc3-inventory.txt --format csv --header -service dash
	`
)

func init() {
	f := HostFlags{}
	cmd := &cobra.Command{
		Use:   "hosts <inventory file>",
		Short: "List UC3 hosts (all, or by service)",
		Long: "List UC3 hosts (all, or by service) from inventory file",
		Example: formatHelp(hostsExamples, "  "),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return f.PrintInventory(args[0])
		},
	}
	cmdFlags := cmd.Flags()
	cmdFlags.StringVarP(&f.service, "service", "s", "", "filter to specified service")
	cmdFlags.StringVarP(&f.FormatStr, "format", "f", output.Default.Name(), formatFlagUsage)
	cmdFlags.BoolVar(&f.Header, "header", false, "include header")
	cmdFlags.BoolVar(&f.Footer, "footer", false, "include footer")
	Root().AddCommand(cmd)
}
