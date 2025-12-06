package proxy

import (
	"net"
	"net/http"
	"net/url"
)

type singleTargetDirector struct {
	targetURL *url.URL
}

type director interface {
	createRequest(r *http.Request) (*http.Request, error)
}

// createRequest creates the new request with the target url and headers
func (d *singleTargetDirector) createRequest(r *http.Request) (*http.Request, error) {
	reqURL := d.targetURL.ResolveReference(r.URL)
	httpReq, err := http.NewRequest(r.Method, reqURL.String(), r.Body)
	if err != nil {
		return nil, err
	}
	removeHopByHopHeaders(r.Header)
	for k, v := range r.Header {
		httpReq.Header[k] = v
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	//sets the X-Forwarded-For header with requesters ip
	httpReq.Header.Set("X-Forwarded-For", ip)
	return httpReq, nil
}
