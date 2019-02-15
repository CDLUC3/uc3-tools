package cmd

import (
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
	"github.com/spf13/pflag"
	"strings"
)

type Flags struct {
	FormatStr string
	Header    bool
	Footer    bool
}

func (f *Flags) AddTo(cmdFlags *pflag.FlagSet) {
	cmdFlags.SortFlags = false

	formatFlagUsage := fmt.Sprintf("output format (%v)", strings.Join(output.StandardFormats(), ", "))
	cmdFlags.StringVarP(&f.FormatStr, "format", "f", output.Default.Name(), formatFlagUsage)
	cmdFlags.BoolVar(&f.Header, "header", false, "include header")
	cmdFlags.BoolVar(&f.Footer, "footer", false, "include footer")
}

