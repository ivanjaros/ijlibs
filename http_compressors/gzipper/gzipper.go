// Source: https://gist.github.com/CJEnright/bc2d8b8dc0c1389a9feeddb110f822d7
package gzipper

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const Name = "gzip"

// this pool uses the best compression setting.
// this is because we need to compress a/v data, which are large,
// but each segment is at least 2 seconds long, so we actually do have time
// to use the best possible compression on our end and save the traffic costs.
var gzPool = sync.Pool{
	New: func() interface{} {
		w, _ := gzip.NewWriterLevel(ioutil.Discard, gzip.BestCompression)
		return w
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), Name) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", Name)

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
