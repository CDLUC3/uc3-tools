package shared

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type errorsByMsg []error

func (s errorsByMsg) Len() int           { return len(s) }
func (s errorsByMsg) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s errorsByMsg) Less(i, j int) bool { return strings.Compare(s[i].Error(), s[j].Error()) < 0 }
func (s errorsByMsg) Eq(i, j int) bool   { return s[i].Error() == s[j].Error() }
func (s errorsByMsg) Copy(j, i int)      { s[j] = s[i] }

//noinspection GoUnhandledErrorResult
func PrintErrors(errors []error) {
	Deduplicate(errorsByMsg(errors), func(len int) {
		errors = errors[:len]
	})
	if len(errors) > 0 {
		w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
		fmt.Fprintf(w, "%d errors:\n", len(errors))
		for i, err := range errors {
			fmt.Fprintf(w, "%d. %v\n", i+1, err)
		}
		w.Flush()
	}
}
