package server

import (
	"compress/gzip"
	"net/http"
)

type GzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (gzip *GzipResponseWriter) Write(b []byte) (int, error) {
	return gzip.Writer.Write(b)
}
