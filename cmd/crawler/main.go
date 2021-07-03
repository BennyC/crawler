package main

import (
	"flag"
	"log"
	"os"

	"github.com/bennyc/crawler"
)

func main() {
	pType := flag.String("print", "text", "available printers: json, text")
	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		log.Fatalf("requires url to get started!")
	}

	var p crawler.Printer
	switch *pType {
	case "text":
		p = crawler.TextPrinter{}
	case "json":
		p = crawler.JsonPrinter{}
	default:
		log.Fatalf("invalid printer type: %s", p)
	}

	res := crawler.Crawl(crawler.Options{
		LinkFinder: crawler.SameDomainLinkFinder{},
	}, url)

	p.Print(os.Stdout, res)
}
