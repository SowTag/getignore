package utils

import (
	"fmt"
	"os"
	"text/tabwriter"

	"golang.org/x/term"
)

func PrintDynamicTable(items []string) {
	const padding = 2
	const assumedFallbackTerminalWidth = 120

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = assumedFallbackTerminalWidth
	}

	longestWordLength := 0
	for _, item := range items {
		if len(item) > longestWordLength {
			longestWordLength = len(item)
		}
	}

	colWidth := longestWordLength + padding
	numCols := width / colWidth

	if numCols == 0 {
		numCols = 1
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for i, item := range items {
		fmt.Fprintf(w, "%s\t", item)

		if (i+1)%numCols == 0 {
			fmt.Fprintln(w)
		}
	}

	if len(items)%numCols != 0 {
		fmt.Fprintln(w)
	}
	w.Flush()
}
