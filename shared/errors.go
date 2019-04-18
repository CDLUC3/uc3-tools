package shared

import (
	"fmt"
	"os"
	"text/tabwriter"
)

//noinspection GoUnhandledErrorResult
func PrintErrors(errors []error) {
	if len(errors) > 0 {
		w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
		fmt.Fprintf(w, "%d errors:\n", len(errors))
		for i, err := range errors {
			fmt.Fprintf(w, "%d. %v\n", i+1, err)
		}
		w.Flush()
	}
}
