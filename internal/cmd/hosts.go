package cmd

import (
	"github.com/dmolesUC3/uc3-system-info/internal/hosts"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command {
		Use: "hosts <FILE>",
		Short: "list UC3 hosts",
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inv, err := hosts.NewInventory(args[0])
			if err != nil {
				return err
			}
			inv.Print()
			return nil
		},
	}

	rootCmd.AddCommand(cmd)
}


