package proxy

import (
	"net/http"
	"net/url"
)

func NewProxy(cfg TransportConfig, targetURL *url.URL) http.Handler {
	return &proxy{
		director: &singleTargetDirector{
			targetURL: targetURL,
		},
		transporter: newTransporter(cfg),
	}
}
