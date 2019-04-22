package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/git"
	"github.com/dmolesUC3/mrt-build-info/jenkins"
	. "github.com/dmolesUC3/mrt-build-info/shared"
	"github.com/spf13/cobra"
	"os"
)

// TODO: 'find' command (by job, artifactId, ...?)

var rootCmd = &cobra.Command{
	Use:   "mrt-build-info",
	Short: "Merritt build info",
	Long:  "Tools for gathering Merritt build information",
}

func ServerFrom(args []string) (server jenkins.JenkinsServer, err error) {
	if len(args) == 0 {
		server = jenkins.DefaultServer()
	} else {
		server, err = jenkins.ServerFromUrl(args[0])
	}
	return
}

func AddCommand(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&Flags.Verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&git.FullSHA, "full-sha", "f", false, "don't abbreviate SHA hashes in URLs")
	cmd.Flags().StringVarP(&git.Token, "token", "t", "", "GitHub API token (https://github.com/settings/tokens)")
	cmd.Flags().BoolVar(&Flags.TSV, "tsv", false, "tab-separated output (default is fixed-width)")

	rootCmd.AddCommand(cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
