package transport

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/n-creativesystem/go-fwncs"
)

type dumpTransport struct {
	log       fwncs.ILogger
	Transport http.RoundTripper
}

func (t *dumpTransport) transport() http.RoundTripper {
	if t.Transport == nil {
		return http.DefaultTransport
	}
	return t.Transport
}

func (t *dumpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.log.Info(fmt.Sprintf("Connected to %v", req.URL))
	dump := func(b []byte) {
		dumps := strings.Split(string(b), "\n")
		for _, dump := range dumps {
			t.log.Debug(dump)
		}
	}
	// リクエストの送信内容を表示
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	dump(b)
	// 実際のリクエストを送信
	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		return nil, err
	}
	// レスポンス内容を表示
	b, err = httputil.DumpResponse(resp, true)
	dump(b)

	return resp, err
}

func NewDumpTransport(log fwncs.ILogger, transport http.RoundTripper) http.RoundTripper {
	return &dumpTransport{
		Transport: transport,
		log:       log,
	}
}
