package internal

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-plugin"
	"github.com/oidc-proxy-ecosystem/proxy-server/plugins"
	"github.com/oidc-proxy-ecosystem/proxy-server/shared"
	"google.golang.org/grpc"
)

type ResponsePlugin struct {
	plugin.Plugin
	Impl shared.Response
}

var _ plugin.GRPCPlugin = (*ResponsePlugin)(nil)

func (p *ResponsePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	plugins.RegisterResponseServer(s, &ResponseGRPCServer{Impl: p.Impl})
	return nil
}

func (p *ResponsePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ResponseGRPCClient{client: plugins.NewResponseClient(c)}, nil
}

type ResponseGRPCServer struct {
	Impl shared.Response
	plugins.UnimplementedResponseServer
}

var _ plugins.ResponseServer = (*ResponseGRPCServer)(nil)

func (m *ResponseGRPCServer) Modify(ctx context.Context, in *plugins.Input) (*plugins.Output, error) {
	header := make(map[string][]string, len(in.Header))
	for key, value := range in.Header {
		header[key] = make([]string, len(value.Value))
		copy(header[key], value.Value)
	}
	rBody := gzipRead(in.Body)
	r := m.Impl.Modify(in.URL, in.Method, header, rBody)
	replyHeader := make(map[string]*plugins.Values, len(r.Header))
	for key, values := range r.Header {
		v := &plugins.Values{
			Value: make([]string, len(values)),
		}
		copy(v.Value, values)
		replyHeader[key] = v
	}
	reply := &plugins.Output{
		Header: replyHeader,
		Body:   gzipWrite(r.Body),
	}
	return reply, nil
}

type ResponseClientInput struct {
	URL    *url.URL
	Method string
	Header http.Header
	Body   []byte
}

type ResponseGRPCClient struct {
	client plugins.ResponseClient
}

func gzipWrite(buf []byte) []byte {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	writer.Write(buf)
	writer.Close()
	return buffer.Bytes()
}

func gzipRead(gzipBuf []byte) []byte {
	buffer := bytes.NewBuffer(gzipBuf)
	reader, _ := gzip.NewReader(buffer)
	output := bytes.Buffer{}
	output.ReadFrom(reader)
	return output.Bytes()
}

func (m *ResponseGRPCClient) Modify(in *ResponseClientInput) (rHeader http.Header, rBody []byte, rErr error) {
	defer func() {
		if err := recover(); err != nil {
			rErr = fmt.Errorf("panic: %v", err)
		}
		rHeader = in.Header.Clone()

	}()
	rHeader = in.Header.Clone()
	rBody = make([]byte, len(in.Body))
	copy(rBody, in.Body)
	rErr = nil
	infHeader := make(map[string]*plugins.Values, len(in.Header))
	for key, values := range in.Header {
		v := &plugins.Values{
			Value: make([]string, len(values)),
		}
		copy(v.Value, values)
		infHeader[key] = v
	}
	inf := &plugins.Input{
		URL:    in.URL.String(),
		Method: in.Method,
		Header: infHeader,
		Body:   gzipWrite(in.Body),
	}
	r, err := m.client.Modify(context.Background(), inf)
	if err != nil {
		rErr = err
		return
	}
	rHeader = make(http.Header, len(r.Header))
	for key, value := range r.Header {
		rHeader[key] = make([]string, len(value.Value))
		copy(rHeader[key], value.Value)
	}
	rBody = gzipRead(r.Body)
	return
}
