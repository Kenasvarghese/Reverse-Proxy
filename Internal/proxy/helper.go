package proxy

import (
	"net/http"
	"strings"
)

var hopByHopHeaders = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

// removeHopByHopHeaders removes the hop by hop headers
func removeHopByHopHeaders(header http.Header) {
	conn := header.Get("Connection")
	for _, token := range strings.Split(conn, ",") {
		token = strings.TrimSpace(token)
		if token != "" {
			header.Del(token)
		}
	}
	for _, h := range hopByHopHeaders {
		header.Del(h)
	}
}
