package crawler

import (
	"sync"
)

// Options for a Crawl
type Options struct {
	// The LinkFinder the Crawler will use to search URLs
	// This could be an implementation for HTML, JSON:API structure with links,
	// Markdown or anything containing links!
	LinkFinder
}

// Crawl return struct
type URLMap struct {
	URLs map[string]DocumentResult

	_mu *sync.Mutex
	_wg *sync.WaitGroup
}

// Check if a url has already been registered within
// the document map
func (d *URLMap) isVisited(url string) bool {
	d._mu.Lock()
	defer d._mu.Unlock()

	if _, ok := d.URLs[url]; !ok {
		return false
	}

	return true
}

// Store the DocumentResult for a given url
func (d *URLMap) store(url string, res DocumentResult) {
	d._mu.Lock()
	defer d._mu.Unlock()

	d.URLs[url] = res
}

// Crawl will recursively crawl all links found from the baseurl
func Crawl(o Options, ch chan DocumentResult, baseurl string) *URLMap {
	dm := &URLMap{
		URLs: map[string]DocumentResult{},

		_mu: &sync.Mutex{},
		_wg: &sync.WaitGroup{},
	}

	dm._wg.Add(1)
	go crawl(o, ch, dm, baseurl)
	dm._wg.Wait()

	close(ch)

	return dm
}

// Recursion function for the exported Crawl
func crawl(o Options, ch chan DocumentResult, dm *URLMap, url string) {
	defer dm._wg.Done()

	res, _ := o.LinkFinder.LinksTo(url)
	dm.store(url, res)

	ch <- res

	for _, u := range res.FoundLinks {
		if !dm.isVisited(u.RawURL) {
			dm._wg.Add(1)
			go crawl(o, ch, dm, u.RawURL)
		}
	}
}
