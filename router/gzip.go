package router

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/n-creativesystem/go-fwncs"
	"github.com/n-creativesystem/go-fwncs/constant"
)

type gzipWriter struct {
	io.Writer
	fwncs.ResponseWriter
}

func (w *gzipWriter) WriteHeader(code int) {
	if code == http.StatusNoContent {
		w.ResponseWriter.Header().Del(constant.HeaderContentEncoding)
	}
	w.Header().Del(constant.HeaderContentLength)
	w.ResponseWriter.WriteHeader(code)
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	if w.Header().Get(constant.HeaderContentType) == "" {
		w.Header().Set(constant.HeaderContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

const gzipScheme = "gzip"

func GzipMiddleware() fwncs.HandlerFunc {
	return func(c fwncs.Context) {
		if strings.Contains(c.Header().Get(constant.HeaderAcceptEncoding), gzipScheme) {
			c.SetHeader(constant.HeaderContentEncoding, gzipScheme)
			gz := gzip.NewWriter(c.Writer())
			defer gz.Close()
			gzw := &gzipWriter{
				Writer:         gz,
				ResponseWriter: c.Writer(),
			}
			c.SetWriter(gzw)
			c.Next()
		} else {
			c.Next()
		}
	}
}
