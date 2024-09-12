package src

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/algrvvv/go-sandbox/src/logger"
)

type wrappedResponseWriter struct {
	statusCode int
	http.ResponseWriter
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *wrappedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("failed to implements hijacker")
	}
	return hijacker.Hijack()
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		logger.Infof("%d %s %s %s", wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func GetNewServer() *http.Server {
	serv := http.NewServeMux()

	serv.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	serv.HandleFunc("GET /", indexHandler)
	serv.HandleFunc("GET /offline", offlineHandler)
	serv.HandleFunc("POST /offline/run", offlineRunHandler)
	serv.HandleFunc("GET /new", newOnlineHandler)
	serv.HandleFunc("GET /online", onlineHandler)
	serv.HandleFunc("POST /online/run", onlineRunHandler)
	serv.HandleFunc("GET /connect", connectOnlineHandler)
	serv.HandleFunc("GET /ws", wsHandler)

	serv.HandleFunc("POST /gofmt", goFmtHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(serv),
	}

	return server
}
