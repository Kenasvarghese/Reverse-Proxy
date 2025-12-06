package proxy

import (
	"net/http"
	"net/url"
)

func NewProxy(targetURL *url.URL) http.Handler {
	return &proxy{
		director: &singleTargetDirector{
			targetURL: targetURL,
		},
		transporter: http.DefaultClient,
	}
}
