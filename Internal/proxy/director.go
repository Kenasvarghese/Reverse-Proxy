package proxy

import (
	"net"
	"net/http"
	"net/url"
)

type director interface {
	createRequest(r *http.Request) (*http.Request, error)
}

type singleTargetDirector struct {
	targetURL *url.URL
}

// createRequest creates the new request with the target url and headers
func (d *singleTargetDirector) createRequest(r *http.Request) (*http.Request, error) {
	reqURL := d.targetURL.ResolveReference(r.URL)
	httpReq, err := http.NewRequestWithContext(r.Context(), r.Method, reqURL.String(), r.Body)
	if err != nil {
		return nil, err
	}
	for k, v := range r.Header {
		httpReq.Header[k] = v
	}
	removeHopByHopHeaders(httpReq.Header)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = ""
	}
	//sets the X-Forwarded-For header with requesters ip
	httpReq.Header.Set("X-Forwarded-For", ip)
	return httpReq, nil
}
