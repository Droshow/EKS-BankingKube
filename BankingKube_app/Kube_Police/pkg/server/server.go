package server

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// LoggingMiddleware logs incoming HTTP requests and their response details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Capture the response status code and size
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		log.Printf("Request: Method=%s, URL=%s, Status=%d, Duration=%s",
			r.Method, r.URL.Path, rw.statusCode, time.Since(startTime))
	})
}

// responseWriter is a custom http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code of the response
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// NewServer initializes and returns a new HTTP server with TLS configuration and logging
func NewServer(certFile, keyFile string) *http.Server {
	mux := http.NewServeMux()

	// Register health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Webhook server is healthy"))
	})

	// Add LoggingMiddleware to all requests
	loggedMux := LoggingMiddleware(mux)

	return &http.Server{
		Addr:    ":8443",
		Handler: loggedMux, // Wrap mux with logging middleware
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
}

// StartServer starts the HTTP server and implements graceful shutdown
func StartServer(server *http.Server) {
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()
	log.Printf("Server is listening on %s", server.Addr)

	// Graceful shutdown on SIGINT
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
