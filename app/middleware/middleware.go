package middleware

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Login struct {
	Logpath io.Writer
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

type GzipMiddleware struct {
	Next http.Handler
}

type TimeoutMiddleware struct {
	Next http.Handler
}

var DefaultTimeoutHandler = http.HandlerFunc(
	func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusGatewayTimeout)
		res.Write([]byte("Service timeout"))
	})

type timeoutWriter struct {
	rw http.ResponseWriter

	status int
	buf    *bytes.Buffer
}

func JsonHandler(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
		return
	})
}

func (gm *GzipMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if gm.Next == nil {
		gm.Next = http.DefaultServeMux
	}

	encodings := r.Header.Get("Accept-Encoding")
	if !strings.Contains(encodings, "gzip") {
		gm.Next.ServeHTTP(w, r)
		return
	}
	w.Header().Add("Content-Encoding", "gzip")
	gzipWriter := gzip.NewWriter(w)
	defer gzipWriter.Close()
	grw := gzipResponseWriter{
		ResponseWriter: w,
		Writer:         gzipWriter,
	}
	gm.Next.ServeHTTP(grw, r)
}

func (grw gzipResponseWriter) Write(data []byte) (int, error) {
	return grw.Writer.Write(data)
}

func (tm TimeoutMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if tm.Next == nil {
		tm.Next = http.DefaultServeMux
	}

	ctx := r.Context()
	ctx, _ = context.WithTimeout(ctx, 30*time.Second)
	r.WithContext(ctx)
	ch := make(chan struct{})
	go func() {
		tm.Next.ServeHTTP(w, r)
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
	}
}

func TimeoutHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tCtx, tCancel := context.WithTimeout(ctx, 30*time.Second)
		cCtx, cCancel := context.WithCancel(ctx)
		r.WithContext(ctx)

		defer tCancel()
		//tw := &timeoutWriter{rw: w, buf: bytes.NewBuffer(nil)}

		go func() {
			next.ServeHTTP(w, r)
			cCancel()
		}()
		select {
		case <-cCtx.Done():
			/*w.WriteHeader(tw.status)
			w.Write(tw.buf.Bytes())*/
                        return
		case <-tCtx.Done():
			if err := tCtx.Err(); err == context.DeadlineExceeded {
				cCancel()
				w.WriteHeader(http.StatusGatewayTimeout)
				w.Write([]byte("Request Timeout, it is taking more than anticipated time to serve request\n"))
			}
		}
	})
}

func GzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		encodings := r.Header.Get("Accept-Encoding")
		if !strings.Contains(encodings, "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Add("Content-Encoding", "gzip")
		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()
		grw := gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gzipWriter,
		}
		next.ServeHTTP(grw, r)
	})
}

func (loger *Login) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		next.ServeHTTP(w, r)

		newlog := log.New(loger.Logpath, " [INFO ]", 0)
		newlog.SetPrefix(logformatter(r, start))
		newlog.Println()
		return
	})
}

func logformatter(r *http.Request, start time.Time) string {
	return "[" + start.Format("2006-01-02 15:04:05") + "]" + " - " + r.Method + " " + r.RequestURI + " " + r.RemoteAddr + " " + time.Since(start).String() + " " + r.Proto 
}
