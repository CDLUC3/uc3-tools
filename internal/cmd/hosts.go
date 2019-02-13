package cmd

import (
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/outputfmt"
	"github.com/dmolesUC3/uc3-system-info/internal/hosts"
	"github.com/spf13/cobra"
	"strings"
)

type Hosts struct {
	formatStr string
	header    bool
	footer bool
	service string
}

func (h *Hosts) PrintHosts(invPath string) error {
	format, err := outputfmt.ToFormat(h.formatStr)
	if err != nil {
		return err
	}
	inv, err := hosts.NewInventory(invPath)
	if err != nil {
		return err
	}
	inv.Print(format, h.header, h.footer, h.service)
	return nil
}

func init() {
	var h Hosts
	cmd := &cobra.Command{
		Use:   "hosts <FILE>",
		Short: "list UC3 hosts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.PrintHosts(args[0])
		},
	}
	formatFlagUsage := fmt.Sprintf("output format (%v)", strings.Join(outputfmt.StandardFormats(), ", "))
	cmd.Flags().StringVarP(&h.formatStr, "format", "f", outputfmt.Default.Name(), formatFlagUsage)
	cmd.Flags().BoolVar(&h.header, "header", false, "include header")
	cmd.Flags().BoolVar(&h.footer, "footer", false, "include footer")
	cmd.Flags().StringVarP(&h.service, "service", "s", "", "filter to specified service")
	rootCmd.AddCommand(cmd)
}
