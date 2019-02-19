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

// Login will hold the writer and will be used while writing logs.
type Login struct {
	Logpath io.Writer
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

// GzipMiddleware holds the http handler which has to be compressed.
type GzipMiddleware struct {
	Next http.Handler
}

// TimeoutMiddleware holds the http handler for which timeout has to be set.
type TimeoutMiddleware struct {
	Next http.Handler
}

var (
	// DefaultTimeoutHandler holds the message that will be thrown when processing time cross the timeout.
	DefaultTimeoutHandler = http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusGatewayTimeout)
			res.Write([]byte("Service timeout"))
		})
)

type timeoutWriter struct {
	rw http.ResponseWriter

	status int
	buf    *bytes.Buffer
}

// JsonHandler will set the header with content-type json.
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

// TimeoutHandler will be responsible for timing out the request if it takes annoying time to serve.
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

// GzipHandler will help in compressing the response served as per the call made to neuron.
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

// Logger will log all the request served by the neuron through API.
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
