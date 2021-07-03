package crawler

import "errors"

var (
	ErrURLParse     = errors.New("unable to parse url")
	ErrHTTPHasIssue = errors.New("url did not return ok status or incorrect content type")
)
