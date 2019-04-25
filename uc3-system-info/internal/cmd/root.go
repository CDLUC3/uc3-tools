package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

var rootCmd *cobra.Command

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Root() *cobra.Command {
	if rootCmd == nil {
		rc :=  &cobra.Command{
			Use: "uc3-system-info",
			Short: "uc3-system-info: a tool for generating UC3 system info reports",
		}
		rc.Flags().SortFlags = false // TODO: figure out why this isn't respected
		rootCmd = rc
	}
	return rootCmd
}

func formatHelp(text, indent string) string {
	formattedText := regexp.MustCompile(`(?m)^[\t ]+`).ReplaceAllString(text, indent)
	formattedText = strings.TrimPrefix(formattedText, "\n")
	formattedText = strings.TrimSuffix(formattedText, "\n")
	return formattedText
}
