package cmd

import (
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
	"github.com/dmolesUC3/uc3-system-info/internal/hosts"
	"github.com/spf13/cobra"
	"strings"
)

type HostFlags struct {
	formatStr string
	header    bool
	footer bool
	service string
}

func (f *HostFlags) PrintInventory(invPath string) error {
	format, err := output.ToFormat(f.formatStr)
	if err != nil {
		return err
	}
	inv, err := hosts.NewInventory(invPath)
	if err != nil {
		return err
	}
	return inv.Print(format, f.header, f.footer, f.service)
}

func init() {
	var h HostFlags
	cmd := &cobra.Command{
		Use:   "hosts <inventory file>",
		Short: "List UC3 hosts",
		Long: "List UC3 hosts from inventory file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.PrintInventory(args[0])
		},
	}
	cmdFlags := cmd.Flags()
	cmdFlags.SortFlags = false

	formatFlagUsage := fmt.Sprintf("output format (%v)", strings.Join(output.StandardFormats(), ", "))
	cmdFlags.StringVarP(&h.formatStr, "format", "f", output.Default.Name(), formatFlagUsage)
	cmdFlags.StringVarP(&h.service, "service", "s", "", "filter to specified service")
	cmdFlags.BoolVar(&h.header, "header", false, "include header")
	cmdFlags.BoolVar(&h.footer, "footer", false, "include footer")
	rootCmd.AddCommand(cmd)
}
