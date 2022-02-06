// Source: https://gist.github.com/CJEnright/bc2d8b8dc0c1389a9feeddb110f822d7
package deflater

import (
	"github.com/klauspost/compress/zlib"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const Name = "deflate"

// this pool uses the best compression setting.
// this is because we need to compress a/v data, which are large,
// but each segment is at least 2 seconds long, so we actually do have time
// to use the best possible compression on our end and save the traffic costs.
var brPool = sync.Pool{
	New: func() interface{} {
		w, _ := zlib.NewWriterLevel(ioutil.Discard, zlib.BestCompression)
		return w
	},
}

type defResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *defResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *defResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *defResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func Deflate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), Name) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", Name)

		def := brPool.Get().(*zlib.Writer)
		defer brPool.Put(def)

		def.Reset(w)
		defer def.Close()

		next.ServeHTTP(&defResponseWriter{ResponseWriter: w, Writer: def}, r)
	})
}
