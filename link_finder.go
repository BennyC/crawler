package crawler

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type LinkFinder interface {
	// LinksTo will find all links within a document body of the given url
	// and return them in a LinkResult
	//
	// LinkFinder can be opinionated on what links it will find
	LinksTo(string) (DocumentResult, error)
}

type Link struct {
	RawURL string
}

type DocumentResult struct {
	// Did the document have issues with parsing or fetching
	HasError bool

	// Any links which are found within a document, which the LinkFinder
	// determines relevant
	//
	// SameDomainLinkFinder will only return URLs for the same domain
	FoundLinks []Link
}

type SameDomainLinkFinder struct{}

// SameDomainLinkFinder will only return links that match the domain
// it was asked to search
func (finder SameDomainLinkFinder) LinksTo(rawurl string) (DocumentResult, error) {
	// validate and build a url
	base, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return DocumentResult{HasError: true}, ErrURLParse
	}

	// fetch the document and build a HTML node from it
	// allowing us to xpath the results to find all links within the document
	doc, err := finder.fetchNode(base)
	if err != nil {
		return DocumentResult{HasError: true}, err
	}

	var links []Link
	nodes := htmlquery.Find(doc, "//a[@href]")

	for _, n := range nodes {
		// if there is an invalid link in the document
		// ignore it, and continue to next node
		l, err := url.ParseRequestURI(htmlquery.SelectAttr(n, "href"))
		if err != nil {
			continue
		}

		// resolve the url against the url that was initially searched for
		// to determine if it is the same host/domain
		resolvedURL := base.ResolveReference(l)
		if resolvedURL.Host == base.Host {
			links = append(links, Link{
				RawURL: resolvedURL.String(),
			})
		}
	}

	return DocumentResult{
		HasError:   false,
		FoundLinks: links,
	}, nil
}

// fetchNode will GET a given url and use the response body to create a *html.Node
//
// An error will be returned for:
//     - any http failures
//     - any document that returns a status code outside of the 200 range
//     - any document that does not have a content type of text/html
func (finder SameDomainLinkFinder) fetchNode(baseurl *url.URL) (*html.Node, error) {
	resp, err := http.Get(baseurl.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	statusOk := (resp.StatusCode >= 200 && resp.StatusCode < 300) && strings.Contains(contentType, "text/html")
	if !statusOk {
		return nil, ErrHTTPHasIssue
	}

	return htmlquery.Parse(resp.Body)
}
