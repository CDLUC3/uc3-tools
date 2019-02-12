package cmd

import (
	"github.com/dmolesUC3/mrt-system-info/internal/hosts"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command {
		Use: "hosts",
		Short: "list UC3 hosts",
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return hosts.NewInventory(args[0]).Print()
		},
	}

	rootCmd.AddCommand(cmd)
}


