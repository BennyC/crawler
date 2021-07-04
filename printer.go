package crawler

import (
	"encoding/json"
	"fmt"
)

type Printer interface {
	// Print will output the DocumentResult in a given format to
	// an io.Writer
	//
	// Some printers such as the JsonPrinter, may store the output
	// until it is ready/complete to send
	Print(chan DocumentResult)
}

type TextPrinter struct{}

// TextPrinter will output the URLMap in a simple format
// listing all of the URLs visited and what links were found on
// that URL
func (t TextPrinter) Print(c chan DocumentResult) {
	for u := range c {
		// Did the URL have issues when visited? Wasn't a HTML document
		// or it didn't load
		// \u2713 - checkmark
		// \u2613 - cross
		hasErrorChar := "\u2713"
		if u.HasError {
			hasErrorChar = "\u2613"
		}

		fmt.Printf("[%s] %s\n", hasErrorChar, u.URL)
		for _, f := range u.FoundLinks {
			fmt.Printf("   - %s\n", f)
		}
	}
}

type JsonPrinter struct{}

// JsonPrinter will output URLMap exactly as it would be marshalled
// This could open for change, to allow the result set to be returned over
// a http/websocket protocol
func (t JsonPrinter) Print(c chan DocumentResult) {
	var dr []DocumentResult
	for u := range c {
		dr = append(dr, u)
	}

	json, _ := json.MarshalIndent(dr, "", "  ")
	fmt.Print(string(json))
}
