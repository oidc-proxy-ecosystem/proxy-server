package router

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/n-creativesystem/go-fwncs"
	"github.com/n-creativesystem/go-fwncs/constant"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal/plugins"
)

func modifyResponsePlugins(responsePlugins []*plugins.PluginClient) func(r *http.Response) error {
	log := fwncs.DefaultLogger
	return func(r *http.Response) error {
		var reader io.Reader
		if r.Body != nil && strings.Contains(r.Header.Get(constant.HeaderContentEncoding), "gzip") {
			r.Header.Del(constant.HeaderContentEncoding)
			if gr, err := gzip.NewReader(r.Body); err != nil {
				reader = r.Body
			} else {
				defer gr.Close()
				buf := new(bytes.Buffer)
				io.Copy(buf, gr)
				reader = buf
			}
		} else {
			reader = r.Body
		}
		buf := new(bytes.Buffer)
		abuf := io.TeeReader(reader, buf)
		for _, client := range responsePlugins {
			rpcClient := client.GetRpcClient()
			raw, err := rpcClient.Dispense("response_modify")
			if err != nil {
				log.Error(err)
				continue
			}
			protoVer := client.GetClient().NegotiatedVersion()
			switch protoVer {
			case 1:
				b, _ := io.ReadAll(abuf)
				p := raw.(*internal.ResponseGRPCClient)
				rHeader, rBody, err := p.Modify(&internal.ResponseClientInput{
					URL:    r.Request.URL,
					Method: r.Request.Method,
					Header: r.Header.Clone(),
					Body:   b,
				})
				if err == nil {
					r.Header = rHeader.Clone()
					r.Body = io.NopCloser(bytes.NewReader(rBody))
					r.Header.Set(constant.HeaderContentLength, strconv.Itoa(len(rBody)))
				} else {
					log.Warning(err)
				}
			default:
				log.Error(fmt.Sprintf("no support version: %d", protoVer))
			}
		}
		return nil

	}
}
