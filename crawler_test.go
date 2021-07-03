package crawler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bennyc/crawler"
	"github.com/stretchr/testify/assert"
)

// Create a new static file server, using the html files within
// the fixtures dir
func newStaticServer(p string) http.Handler {
	return http.FileServer(http.Dir(p))
}

func TestCrawlWillReturnTheCorrectTreeForSimpleSite(t *testing.T) {
	ts := httptest.NewServer(newStaticServer("./fixtures/simple"))
	defer ts.Close()

	results := crawler.Crawl(crawler.Options{
		LinkFinder: crawler.SameDomainLinkFinder{},
	}, ts.URL+"/a.html")

	assert.Len(t, results.URLs, 5)
}

func TestCrawlWillReturnTheCorrectTreeForSingleSiteWithDeadLinks(t *testing.T) {
	ts := httptest.NewServer(newStaticServer("./fixtures/single"))
	defer ts.Close()

	results := crawler.Crawl(crawler.Options{
		LinkFinder: crawler.SameDomainLinkFinder{},
	}, ts.URL+"/a.html")

	assert.Len(t, results.URLs, 6)
}
