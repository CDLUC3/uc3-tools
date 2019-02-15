package cmd

import (
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
	"github.com/dmolesUC3/uc3-system-info/internal/storage"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
	"time"
)

type NodeFlags struct {
	formatStr string
	header    bool
	footer    bool
}

func (f *NodeFlags) PrintNodes(mrtConfPath string) error {
	format, err := output.ToFormat(f.formatStr)
	if err != nil {
		return err
	}

	mrtConf, err := storage.NewMrtConf(mrtConfPath)
	if err != nil {
		return err
	}

	nodeSets, err := mrtConf.NodeSets()
	if err != nil {
		return err
	}

	for i, nodeSet := range nodeSets {
		fmt.Println(format.SprintTitle(filepath.Base(nodeSet.PropsPath)))
		if f.header {
			fmt.Print(format.SprintHeader("Node number", "Service", "Container"))
		}
		for _, node := range nodeSet.Nodes() {
			fmt.Println(node.Sprint(format))
		}
		if i + 1 < len(nodeSets) {
			fmt.Println()
		}
	}

	if f.footer {
		fmt.Printf("\nGenerated %v\n", time.Now().Format(time.RFC3339))
	}
	return nil
}

func init() {
	var n NodeFlags

	cmd := &cobra.Command{
		Use:   "nodes <path to mrt-conf-prv>",
		Short: "List storage nodes",
		Long: "List storage nodes defined in mrt-conf-prv",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return n.PrintNodes(args[0])
		},
	}
	cmdFlags := cmd.Flags()
	cmdFlags.SortFlags = false

	formatFlagUsage := fmt.Sprintf("output format (%v)", strings.Join(output.StandardFormats(), ", "))
	cmdFlags.StringVarP(&n.formatStr, "format", "f", output.Default.Name(), formatFlagUsage)
	cmdFlags.BoolVar(&n.header, "header", false, "include header")
	cmdFlags.BoolVar(&n.footer, "footer", false, "include footer")
	rootCmd.AddCommand(cmd)
}
