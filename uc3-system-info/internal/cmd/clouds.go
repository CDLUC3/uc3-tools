package cmd

import (
	"fmt"
	"github.com/CDLUC3/uc3-tools/uc3-system-info/internal/output"
	. "github.com/CDLUC3/uc3-tools/uc3-system-info/internal/storage"
	"github.com/spf13/cobra"
	"reflect"
	"sort"
	"time"
)

type CloudFlags struct {
	Flags
}

func (f *CloudFlags) PrintServices() error {
	mrtConfPath := f.ConfPath
	conf, err := NewMrtConf(mrtConfPath)
	if err != nil {
		return err
	}

	format, err := output.ToFormat(f.FormatStr)
	if err != nil {
		return err
	}

	nodeSets, err := conf.NodeSets()
	if err != nil {
		return err
	}

	allServices := map[string]CloudService{}
	for _, nodeSet := range nodeSets {
		services := nodeSet.Services()
		for name, svcP := range services {
			svc := *svcP
			if existing, exists := allServices[name]; exists {
				if !reflect.DeepEqual(existing, svc) {
					return fmt.Errorf("incompatible definitions for service %v:\n\t%v\n\t%v",
						name, existing.Sprint(output.CSV), svcP.Sprint(output.CSV))
				}
			} else {
				allServices[name] = svc
			}
		}
	}

	var allNames []string
	for name := range allServices {
		allNames = append(allNames, name)
	}
	sort.Strings(allNames)

	if f.Header {
		headerStr := format.SprintHeader("Name", "Service type", "Access mode", "Endpoint", "Key", "Secret")
		fmt.Print(headerStr)
	}
	for _, name := range allNames {
		svc := allServices[name]
		svcStr := svc.Sprint(format)
		fmt.Println(svcStr)
	}
	if f.Footer {
		fmt.Printf("\nGenerated %v\n", time.Now().Format(time.RFC3339))
	}
	return nil
}

func init() {
	f := CloudFlags{}

	cmd := &cobra.Command{
		Use:   "clouds",
		Short: "List Merritt cloud storage services",
		Long: "List cloud storage services defined in mrt-conf-prv",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return f.PrintServices()
		},
	}
	cmdFlags := cmd.Flags()
	f.AddTo(cmdFlags)
	Root().AddCommand(cmd)
}
