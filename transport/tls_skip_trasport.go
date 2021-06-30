package transport

import (
	"crypto/tls"
	"net/http"
)

type tlsSkipTransport struct {
	tr http.RoundTripper
}

func (t *tlsSkipTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	tr := t.tr
	if tr == nil {
		if transport, ok := http.DefaultTransport.(*http.Transport); ok {
			if transport.TLSClientConfig == nil {
				transport.TLSClientConfig = &tls.Config{}
			}
			transport.TLSClientConfig.InsecureSkipVerify = true
			tr = transport
		}
	}
	resp, err = tr.RoundTrip(req)
	return
}

func NewTLSSkipVerifyTransport(tr http.RoundTripper) http.RoundTripper {
	return &tlsSkipTransport{
		tr: tr,
	}
}
