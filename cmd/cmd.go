package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/git"
	. "github.com/dmolesUC3/mrt-build-info/shared"
	"github.com/spf13/cobra"
	"os"
)

// TODO: 'find' command (by job, artifactId, ...?)

const valueUnknown = "(unknown)"

var rootCmd = &cobra.Command{
	Use:   "mrt-build-info",
	Short: "Merritt build info",
	Long:  "Tools for gathering Merritt build information",
}

func AddCommand(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&Flags.Job, "job", "j", "", "show info only for specified job")
	cmd.Flags().BoolVarP(&Flags.Verbose, "verbose", "v", false, "verbose output")

	cmd.Flags().BoolVar(&Flags.TSV, "tsv", false, "tab-separated output (default is fixed-width)")

	cmd.Flags().BoolVarP(&git.FullSHA, "full-sha", "f", false, "don't abbreviate SHA hashes in URLs")
	cmd.Flags().StringVarP(&git.Token, "token", "t", "", "GitHub API token (https://github.com/settings/tokens)")

	rootCmd.AddCommand(cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
