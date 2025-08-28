package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			h.ServeHTTP(res, req)
			return
		}

		ow := res

		reqAcceptEncoding := req.Header.Get("Accept-Encoding")
		checkSupportsGzip := strings.Contains(reqAcceptEncoding, "gzip")
		if checkSupportsGzip {
			ow := newCompressWriter(res)
			ow.Header().Set("Content-Encoding", "gzip")
			
			defer ow.Close()
		}

		reqContentEncoding := req.Header.Get("Content-Encoding")
		checkSendsGzip := strings.Contains(reqContentEncoding, "gzip")

		reqContentType := req.Header.Get("Content-Type")
		checkSupportsType := strings.Contains(reqContentType, "application/json") || strings.Contains(reqContentType, "text/plain")
	
		if checkSendsGzip && checkSupportsType {
			cr, err := newCompressReader(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, req)
	})
}
