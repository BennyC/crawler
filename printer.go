package crawler

import (
	"encoding/json"
	"fmt"
	"io"
)

type Printer interface {
	// Print will output the URLMap in a given format to
	// an io.Writer
	Print(io.Writer, *URLMap)
}

type TextPrinter struct{}

// TextPrinter will output the URLMap in a simple format
// listing all of the URLs visited and what links were found on
// that URL
func (t TextPrinter) Print(w io.Writer, d *URLMap) {
	for k, u := range d.URLs {
		// Did the URL have issues when visited? Wasn't a HTML document
		// or it didn't load
		// \u2713 - checkmark
		// \u2613 - cross
		hasErrorChar := "\u2713"
		if u.HasError {
			hasErrorChar = "\u2613"
		}

		fmt.Fprintf(w, "[%s] %s\n", hasErrorChar, k)
		for _, f := range u.FoundLinks {
			fmt.Fprintf(w, "   - %s\n", f)
		}
	}
}

type JsonPrinter struct{}

// JsonPrinter will output URLMap exactly as it would be marshalled
// This could open for change, to allow the result set to be returned over
// a http/websocket protocol
func (t JsonPrinter) Print(w io.Writer, d *URLMap) {
	json, _ := json.MarshalIndent(d, "", "  ")
	fmt.Fprint(w, string(json))
}
