package main

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal"
	"github.com/oidc-proxy-ecosystem/proxy-server/internal/plugins"
	"github.com/oidc-proxy-ecosystem/proxy-server/plugin/modify/files"
	"github.com/oidc-proxy-ecosystem/proxy-server/shared"
)

type Hello struct {
}

var _ shared.Response = (*Hello)(nil)

func (h *Hello) Modify(URL string, method string, header map[string][]string, body []byte) shared.ResponseResult {
	var resultBody []byte
	buf := new(bytes.Buffer)
	httpHeader := http.Header(header)
	if strings.HasPrefix(httpHeader.Get("Content-Type"), "text/html") {
		t := template.Must(template.ParseFS(files.File, "templates/test.html"))
		t.Execute(buf, map[string]string{
			"Username": "admin",
			"Password": "password",
		})
		str := string(body) + buf.String()
		resultBody = []byte(str)
	} else {
		resultBody = body
	}
	return shared.ResponseResult{
		Header: header,
		Body:   resultBody,
	}
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: {
				"response_modify": &internal.ResponsePlugin{Impl: &Hello{}},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
