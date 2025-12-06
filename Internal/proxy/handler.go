package proxy

import (
	"io"
	"log"
	"net/http"
)

type proxy struct {
	director    director
	transporter http.RoundTripper
}

// ServeHTTP implements the http.handler interface
// proxy which forwards the requests to configured url
func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	req, err := p.director.createRequest(r)
	resp, err := p.transporter.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	removeHopByHopHeaders(resp.Header)
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("error copying response body: %v", err)
	}
}
