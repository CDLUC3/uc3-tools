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

	examples := []string{
		"uc3-system-info locate -c ~/Work/mrt-conf-prv -e stg -n 5001 -a ark:/b5072/fk2wq01k85",
		"uc3-system-info locate -c ~/Work/mrt-conf-prv -e stg -n 5001 -a ark:/b5072/fk2wq01k85 -v 1 producer/20151-semestre.csv",
		"uc3-system-info locate -c ~/Work/mrt-conf-prv -e stg -n 9001 -a ark:/99999/fk4kw5kc1z",
		"uc3-system-info locate -c ~/Work/mrt-conf-prv -e stg -n 9001 -a ark:/99999/fk4kw5kc1z -v 1 producer/6GBZeroFile.txt",
	}

	longDesc := []string {
		"Locate a file in Merritt cloud storage based on the service, node number, and object ARK.",
		"(Note that this does not guarantee that the file exists, but only provides the information",
		"necessary to find it.)",
	}

	f := LocateFlags{}
	cmd := &cobra.Command{
		Use:     "locate [filepath]",
		Short:   "Locate object file in Merritt cloud storage",
		Long:    strings.Join(longDesc, "\n"),
		Args:    cobra.MaximumNArgs(1),
		Example: strings.Join(examples, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return f.PrintFileLocation(args[0])
			}
			return f.PrintObjectLocation()
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

func (f *LocateFlags) PrintObjectLocation() error {
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
	serviceDesc := svc.Sprint(output.CSV)

	container, err := node.ContainerFor(f.Ark)
	if err != nil {
		return err
	}

	cliExample, err := node.CLIExampleObject(f.Ark)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, '\t', 0)
	_, _ = fmt.Fprintf(w, "%v:\t%v\n", "Service", serviceDesc)
	_, _ = fmt.Fprintf(w, "%v:\t%v\n", "Container", container)
	_, _ = fmt.Fprintf(w, "%v:\t%v\n", "Example", cliExample)

	return w.Flush()
}

func (f *LocateFlags) PrintFileLocation(filepath string) error {
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
	serviceDesc := svc.Sprint(output.CSV)

	container, err := node.ContainerFor(f.Ark)
	if err != nil {
		return err
	}
	key := node.KeyFor(f.Ark, f.Version, filepath)

	cliExample, err := node.CLIExampleFile(f.Ark, f.Version, filepath)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, '\t', 0)
	_, _ = fmt.Fprintf(w, "%v:\t%v\n", "Service", serviceDesc)
	_, _ = fmt.Fprintf(w, "%v:\t%v\n", "Container", container)
	_, _ = fmt.Fprintf(w, "%v:\t%v\n", "Key", key)
	_, _ = fmt.Fprintf(w, "%v:\t%v\n", "Example", cliExample)

	return w.Flush()
}
