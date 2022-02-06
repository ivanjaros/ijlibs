package http_compressors

import (
	"github.com/ivanjaros/ijlibs/http_compressors/brotler"
	"github.com/ivanjaros/ijlibs/http_compressors/deflater"
	"github.com/ivanjaros/ijlibs/http_compressors/gzipper"
	"net/http"
	"strings"
)

func HttpEncodingHandler(next http.Handler, requireEncoding bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") == "websocket" {
			next.ServeHTTP(w, r)
			return
		}

		type compressor struct {
			Name    string
			Handler func(http.Handler) http.Handler
		}

		// compressors are ordered in case the client accepts multiple formats(most likely)
		compressors := []compressor{
			{Name: brotler.Name, Handler: brotler.Brotli},
			{Name: gzipper.Name, Handler: gzipper.Gzip},
			{Name: deflater.Name, Handler: deflater.Deflate},
		}

		var c func(http.Handler) http.Handler
		for _, v := range compressors {
			if strings.Contains(r.Header.Get("Accept-Encoding"), v.Name) {
				c = v.Handler
				break
			}
		}

		if c == nil && requireEncoding {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		if c == nil {
			next.ServeHTTP(w, r)
		} else {
			c(next).ServeHTTP(w, r)
		}
	})
}
