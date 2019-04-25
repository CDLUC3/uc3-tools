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

func (f *NodeFlags) PrintNodes() error {
	mrtConfPath := f.ConfPath
	mrtConf, err := NewMrtConf(mrtConfPath)
	if err != nil {
		return err
	}

	format, err := output.ToFormat(f.FormatStr)
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
		if i+1 < len(nodeSets) {
			fmt.Println()
		}
	}

	if f.Footer {
		fmt.Printf("\nGenerated %v\n", time.Now().Format(time.RFC3339))
	}
	return nil
}

const (
	nodesLongDesc = `
		List storage nodes defined in mrt-conf-prv.

		Note that for Swift storage, the container shown is only a base name; to get 
		the actual container for a specific object, this base must be suffixed with 
		the first three digits of the MD5 sum of the object ARK. The locate command
		can be used to determine the actual container for an object.
	`

	nodesExamples = `
		uc3-system-info nodes -c ~/Work/mrt-conf-prv
		uc3-system-info nodes -c ~/Work/mrt-conf-prv --header --footer --format md
	`
)

func init() {
	f := NodeFlags{}

	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "List Merritt storage nodes",
		Long:  formatHelp(nodesLongDesc, ""),
		Example: formatHelp(nodesExamples, "  "),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return f.PrintNodes()
		},
	}
	cmdFlags := cmd.Flags()
	f.AddTo(cmdFlags)
	Root().AddCommand(cmd)
}
