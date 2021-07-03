package crawler_test

import (
	"net/http/httptest"
	"testing"

	"github.com/bennyc/crawler"
	"github.com/stretchr/testify/assert"
)

func TestInvalidURL(t *testing.T) {
	ts := httptest.NewServer(newStaticServer("./fixtures/invalid"))
	defer ts.Close()

	lf := crawler.SameDomainLinkFinder{}
	res, err := lf.LinksTo("http//asd")
	assert.Error(t, err)
	assert.Lenf(t, res.FoundLinks, 0, "same domain links expected: %d, got: %d", 0, len(res.FoundLinks))
}

func TestSameDomainLinkFinderWithErrors(t *testing.T) {
	ts := httptest.NewServer(newStaticServer("./fixtures/invalid"))
	defer ts.Close()

	lf := crawler.SameDomainLinkFinder{}
	res, err := lf.LinksTo(ts.URL + "/invalid.json")
	assert.Error(t, err)
	assert.Lenf(t, res.FoundLinks, 0, "same domain links expected: %d, got: %d", 0, len(res.FoundLinks))
}

func TestSameDomainLinkFinderOnlyFindsLinksForTheSameDomain(t *testing.T) {
	ts := httptest.NewServer(newStaticServer("./fixtures/simple"))
	defer ts.Close()

	testCases := []struct {
		Count int
		URL   string
	}{
		{
			Count: 3,
			URL:   ts.URL + "/a.html",
		},
		{
			Count: 2,
			URL:   ts.URL + "/b.html",
		},
		{
			Count: 2,
			URL:   ts.URL + "/c.html",
		},
		{
			Count: 0,
			URL:   ts.URL + "/d.html",
		},
	}

	lf := crawler.SameDomainLinkFinder{}
	for _, tc := range testCases {
		res, err := lf.LinksTo(tc.URL)

		assert.NoError(t, err)
		assert.Lenf(t, res.FoundLinks, tc.Count, "same domain links expected: %d, got: %d", tc.Count, len(res.FoundLinks))
	}
}
