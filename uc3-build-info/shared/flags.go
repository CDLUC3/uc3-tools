package shared

var Flags struct {
	Job string
	TSV bool
	Short bool
	// TODO: replace explicit Verbose checks with a logger
	Verbose bool
}

const ValueUnknown = "(unknown)"

