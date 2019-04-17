package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// TODO: 'find' command (by job, artifactId, ...?)

var rootCmd = &cobra.Command {
	Use:   "mrt-build-info",
	Short: "Merritt build info",
	Long:  "Tools for gathering Merritt build information",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}