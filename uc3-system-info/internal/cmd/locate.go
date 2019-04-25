package cmd

import (
	"fmt"
	"github.com/CDLUC3/uc3-tools/uc3-system-info/internal/output"
	. "github.com/CDLUC3/uc3-tools/uc3-system-info/internal/storage"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

const (
	locateExamples = `
		uc3-system-info locate -c ~/Work/mrt-conf-prv -n 5001 -a ark:/b5060/d8rp4v
		uc3-system-info locate -c ~/Work/mrt-conf-prv -n 5001 -a ark:/b5060/d8rp4v -v 1 producer/README.txt 
		uc3-system-info locate -c ~/Work/mrt-conf-prv -e stg -n 9001 -s replic -a ark:/99999/fk4kw5kc1z
		uc3-system-info locate -c ~/Work/mrt-conf-prv -e stg -n 9001 -s replic -a ark:/99999/fk4kw5kc1z -v 1 producer/6GBZeroFile.txt
	`
	
	locateLongDesc = `
		Locates an object in Merritt cloud storage based on the service, node
		number, and object ARK, or locate a specific file in that object.

		(Note that this does not guarantee that the object exists, on that
		storage node, but only provides the information necessary to find it.)

		In general, the cloud storage key for a file is of the form
		
		  <ark>|<version>|<file>

		Note that files are stored only under the version in which they were
		originally uploaded and under any versions in which their content was
		changed.

		Note that for Swift containers ending in .__, the base container name
		must be suffixed with the first three digits of the MD5 sum of the object
		ARK to obtain the actual container.

		For objects in S3 and Swift, the locate command will print a sample
		command line for accessing the object, including credentials if
		available.
				
		On macOS, the example command for Swift storage will use the full path
		/usr/local/bin/swift, in order to avoid collisions with the compiler
		for the Swift programming language. If the OpenStack Swift CLI is
		installed somehwere other than /usr/local/bin, the example command may
		need to be modified.
	`
)

func init() {
	f := LocateFlags{}
	cmd := &cobra.Command{
		Use:     "locate [filepath]",
		Short:   "Locate object or file in Merritt cloud storage",
		Long:    formatHelp(locateLongDesc, ""),
		Args:    cobra.MaximumNArgs(1),
		Example: formatHelp(locateExamples, "  "),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return f.PrintFileLocation(args[0])
			}
			return f.PrintObjectLocation()
		},
	}
	cmdFlags := cmd.Flags()
	cmdFlags.SortFlags = false
	cmdFlags.StringVarP(&f.ConfPath, "conf", "c", "", "path to mrt-conf-prv project (required)")
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
