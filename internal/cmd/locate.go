package cmd

import (
	"fmt"
	"github.com/dmolesUC3/uc3-system-info/internal/output"
	. "github.com/dmolesUC3/uc3-system-info/internal/storage"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

func init() {

	// TODO: flags to generate AWS or S3 commands

	examples := []string {
		"uc3-system-info locate producer/Prasad_ucla_0031D_15251.pdf -c ~/Work/mrt-conf-prv -n 9001 -a ark:/99999/fk4qz2hp2t -v 1",
	}

	f := LocateFlags{}
	cmd := &cobra.Command{
		Use:   "locate <filepath>",
		Short: "Locate file in Merritt cloud storage",
		Long:  "Locate a file in Merritt cloud storage based on the service, node number, and object ARK",
		Args:  cobra.ExactArgs(1),
		Example: strings.Join(examples, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return f.PrintLocation(args[0])
		},
	}
	cmdFlags := cmd.Flags()
	cmdFlags.SortFlags = false
	cmdFlags.StringVarP(&f.ConfPath, "conf", "c", "", "path to mrt-conf-prv project")
	cmdFlags.StringVarP(&f.Service, "service", "s", "store", "service (store, replic, audit)")
	cmdFlags.StringVarP(&f.Environment, "env", "e", "prd", "environment (dev, stg, prd)")
	cmdFlags.Int64VarP(&f.NodeNumber, "node", "n", 0, "node number")
	cmdFlags.StringVarP(&f.Ark, "ark", "a", "", "object ARK")
	cmdFlags.IntVarP(&f.Version, "version", "v", 1, "object version")

	Root().AddCommand(cmd)
}

// TODO: pass conf as argument, path as flag?
type LocateFlags struct {
	ConfPath    string
	Service     string
	Environment string

	NodeNumber int64

	Ark     string
	Version int
}

func (f *LocateFlags) PrintLocation(filepath string) error {
	conf, err := NewMrtConf(f.ConfPath)
	if err != nil {
		return err
	}

	node, err := conf.GetNode(f.Environment, f.Service, f.NodeNumber)
	if err != nil {
		return err
	}

	svc := node.Service
	if svc == nil {
		return fmt.Errorf("unable to determine cloud service for node %d", f.NodeNumber)
	}

	container, err := node.ContainerFor(f.Ark)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%v|%d|%v", f.Ark, f.Version, filepath)

	fields := map[string]string{
		"Service":   svc.Sprint(output.CSV),
		"Container": container,
		"Key":       key,
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, '\t', 0)
	for k, v := range fields {
		_, err = fmt.Fprintf(w, "%v:\t%v\n", k, v)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}
