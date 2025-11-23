package proxy

import (
	"io"
	"net/http"
	"net/url"
)

type proxy struct {
	url    *url.URL
	server *http.Server
}

func NewProxyHandler(url *url.URL) http.Handler {
	return &proxy{url: url}
}
func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	newUrl := p.url.ResolveReference(r.URL)
	req, err := http.NewRequest(r.Method, newUrl.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for k, v := range r.Header {
		req.Header[k] = v
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
