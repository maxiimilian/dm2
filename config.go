package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// Regexes for config validation
var reSection *regexp.Regexp = regexp.MustCompile(`\[(.*?)\]`)
var reFilterLine *regexp.Regexp = regexp.MustCompile(`^[+-]\s.*$`)
var reComment *regexp.Regexp = regexp.MustCompile(`^#.*$`)

// Define config modes
const MODE_SKIP int = 0
const MODE_KEEP int = 1

// Read and sanitize config file from `configPath`.
// Depending on the `mode` (SKIP OR KEEP) lines will be skipped or kept if they
// match any of the regular expressions in `reList`
func readSanitizeConfig(configPath string, reList []*regexp.Regexp, mode int) []string {
	InfoLogger.Printf("Reading and sanitizing config file '%s'.", configPath)
	// Create empty config slice
	var config []string

	// Open file and create new scanner to iterate over content
	f, err := os.Open(configPath)
	checkError(err)

	var i int = 0
	s := bufio.NewScanner(f)

	for s.Scan() {
		i++
		line := s.Text()

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		switch mode {
		case MODE_KEEP:
			// Keep ONLY line that match
			if matchStringList(line, reList) {
				config = append(config, line)
			}
			continue
		case MODE_SKIP:
			// Skip lines that match
			if matchStringList(line, reList) {
				continue
			}
			// Keep ALL OTHER lines
			config = append(config, line)
			continue
		}
	}
	return config
}

// Match string `s` against list of regular expressions.
// Returns true on first match or false if no expression matches.
func matchStringList(s string, reList []*regexp.Regexp) bool {
	for _, re := range reList {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}

func parseDatasetConfig(config []string, globalFilters *[]string) []*Dataset {
	datasets := make([]*Dataset, 0)
	var ds *Dataset
	var name string

	for _, line := range config {
		if reSection.MatchString(line) {
			// New dataset
			name = reSection.FindStringSubmatch(line)[1] // first group is name
			ds = &Dataset{
				name:          name,
				filters:       make([]string, 0),
				globalFilters: globalFilters,
			}
			datasets = append(datasets, ds)
			continue
		}

		// Config was already validated during reading, so only filters follow after section
		ds.filters = append(ds.filters, line)
	}

	return datasets
}

// Load datasets from file and return slice of Dataset pointers
func loadDatasetConfig(configPath string, globalFilters *[]string) []*Dataset {
	keep := []*regexp.Regexp{reSection, reFilterLine}
	config := readSanitizeConfig(configPath, keep, MODE_KEEP)
	return parseDatasetConfig(config, globalFilters)
}

// Load list of global ignores from file an return list of strings
func loadGlobalIgnore(ignorePath string) *[]string {
	skip := []*regexp.Regexp{reComment}
	ignores := readSanitizeConfig(ignorePath, skip, MODE_SKIP)

	// Prefix all ignores with '-' to make them valid filter
	for i, line := range ignores {
		ignores[i] = fmt.Sprintf("- %s", line)
	}

	return &ignores
}
