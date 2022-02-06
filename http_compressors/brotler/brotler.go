// Source: https://gist.github.com/CJEnright/bc2d8b8dc0c1389a9feeddb110f822d7
package brotler

import (
	"github.com/andybalholm/brotli"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const Name = "br"

// this pool uses the best compression setting.
// this is because we need to compress a/v data, which are large,
// but each segment is at least 2 seconds long, so we actually do have time
// to use the best possible compression on our end and save the traffic costs.
var brPool = sync.Pool{
	New: func() interface{} {
		return brotli.NewWriterLevel(ioutil.Discard, brotli.BestCompression)
	},
}

type brResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *brResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *brResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *brResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func Brotli(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), Name) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", Name)

		br := brPool.Get().(*brotli.Writer)
		defer brPool.Put(br)

		br.Reset(w)
		defer br.Close()

		next.ServeHTTP(&brResponseWriter{ResponseWriter: w, Writer: br}, r)
	})
}
