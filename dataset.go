package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Dataset struct {
	name          string
	filters       []string
	globalFilters *[]string
}

// WriteTmpFilterFile writes filters of Dataset instance to temporary file and return its path.
func (ds *Dataset) WriteTmpFilterFile() string {
	// Create temp filter list
	var filters []string

	// Global filters do not have to be prefixed. Just append them to the beginning of filter list
	filters = append(filters, *ds.globalFilters...)

	if len(ds.filters) == 0 {
		// Default: No filters -> include everything from dataset directory
		filters = append(filters, fmt.Sprintf("+ %s/**", ds.name))
	} else {
		// Prefix dataset-specific filters with dataset name
		var flag, filter string

		for _, f := range ds.filters {
			flag = string(f[0])
			filter = f[2:]

			filters = append(filters, fmt.Sprintf("%s %s/%s", flag, ds.name, filter))
		}
	}

	// Finally add "exclude rest"
	filters = append(filters, "- **")

	// Open temporary file for writing
	f, err := os.CreateTemp("", "dm_*.txt")
	checkError(err)

	// Write list of filters to file
	w := bufio.NewWriter(f)
	nBytes, err := w.WriteString(strings.Join(filters, "\n"))
	checkError(err)

	err = w.Flush()
	checkError(err)

	DebugLogger.Printf("%d bytes written to `%s`", nBytes, f.Name())

	return f.Name()
}
