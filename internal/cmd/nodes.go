package cmd

import (
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
	. "github.com/dmolesUC3/uc3-system-info/internal/storage"
	"github.com/spf13/cobra"
	"path/filepath"
	"time"
)

type NodeFlags struct {
	Flags
}

func (f *NodeFlags) PrintNodes(mrtConfPath string) error {
	format, err := output.ToFormat(f.FormatStr)
	if err != nil {
		return err
	}

	mrtConf, err := NewMrtConf(mrtConfPath)
	if err != nil {
		return err
	}

	nodeSets, err := mrtConf.NodeSets()
	if err != nil {
		return err
	}

	for i, nodeSet := range nodeSets {
		fmt.Println(format.SprintTitle(filepath.Base(nodeSet.PropsPath)))
		if f.Header {
			fmt.Print(format.SprintHeader("Node number", "Service", "Container"))
		}
		for _, node := range nodeSet.Nodes() {
			fmt.Println(node.Sprint(format))
		}
		if i + 1 < len(nodeSets) {
			fmt.Println()
		}
	}

	if f.Footer {
		fmt.Printf("\nGenerated %v\n", time.Now().Format(time.RFC3339))
	}
	return nil
}

func init() {
	f := NodeFlags{}

	cmd := &cobra.Command{
		Use:   "nodes <path to mrt-conf-prv>",
		Short: "List storage nodes",
		Long: "List storage nodes defined in mrt-conf-prv",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return f.PrintNodes(args[0])
		},
	}
	cmdFlags := cmd.Flags()
	f.AddTo(cmdFlags)
	Root().AddCommand(cmd)
}
