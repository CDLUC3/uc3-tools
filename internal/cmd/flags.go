package cmd

import (
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

type Flags struct {
	ConfPath    string
	FormatStr string
	Header    bool
	Footer    bool
}

var formatFlagUsage = fmt.Sprintf("output format (%v)", strings.Join(output.StandardFormats(), ", "))

func (f *Flags) AddTo(cmdFlags *pflag.FlagSet) {
	cmdFlags.SortFlags = false

	cmdFlags.StringVarP(&f.ConfPath, "conf", "c", "", "path to mrt-conf-prv project (required)")
	_ = cobra.MarkFlagRequired(cmdFlags, "conf")

	cmdFlags.StringVarP(&f.FormatStr, "format", "f", output.Default.Name(), formatFlagUsage)
	cmdFlags.BoolVar(&f.Header, "header", false, "include header")
	cmdFlags.BoolVar(&f.Footer, "footer", false, "include footer")
}

