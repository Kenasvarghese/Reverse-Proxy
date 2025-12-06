package proxy

import "net/http"

type transporter interface {
	Do(req *http.Request) (*http.Response, error)
}
