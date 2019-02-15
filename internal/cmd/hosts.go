package cmd

import (
	. "github.com/dmolesUC3/uc3-system-info/internal/hosts"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
	"github.com/spf13/cobra"
)

type HostFlags struct {
	Flags
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

func init() {
	f := HostFlags{}
	cmd := &cobra.Command{
		Use:   "hosts <inventory file>",
		Short: "List UC3 hosts",
		Long: "List UC3 hosts from inventory file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return f.PrintInventory(args[0])
		},
	}
	cmdFlags := cmd.Flags()
	cmdFlags.StringVarP(&f.service, "service", "s", "", "filter to specified service")
	f.AddTo(cmdFlags)
	Root().AddCommand(cmd)
}
