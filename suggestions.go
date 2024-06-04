package exoskeleton

import (
	"strings"
)

const suggestionsMinimumDistance = 2

// suggestionsFor returns a list of commands that have similar names to the one given.
func (e *Entrypoint) suggestionsFor(typedName string) (suggestions []Command) {
	typedName = strings.Replace(typedName, ":", " ", -1)

	seen := make(map[string]bool)

	// Ignore errors when building up a list of suggestions
	cmds, _ := e.Subcommands()
	for _, cmd := range cmds.Flatten() {
		usage := UsageRelativeTo(cmd, e)
		if seen[usage] {
			continue
		}

		levenshteinDistance := ld(typedName, usage, true)
		suggestByLevenshtein := levenshteinDistance <= suggestionsMinimumDistance
		suggestByPrefix := strings.HasPrefix(strings.ToLower(usage), strings.ToLower(typedName))
		suggestByNamspace := strings.HasSuffix(strings.ToLower(usage), " "+strings.ToLower(typedName))
		if suggestByLevenshtein || suggestByPrefix || suggestByNamspace {
			suggestions = append(suggestions, cmd)
			seen[usage] = true
		}
	}

	return
}

// ld compares two strings and returns the levenshtein distance between them.
func ld(s, t string, ignoreCase bool) int {
	if ignoreCase {
		s = strings.ToLower(s)
		t = strings.ToLower(t)
	}
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(t)]
}
