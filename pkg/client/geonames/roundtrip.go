package geonames

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type loggingRoundTripper struct {
	logger io.Writer
	next   http.RoundTripper
}

func (l loggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	_, err := fmt.Fprintf(l.logger, "[%s] %s %s\n", time.Now().Format(time.ANSIC), r.Method, r.URL)
	if err != nil {
		return nil, err
	}

	return l.next.RoundTrip(r)
}
