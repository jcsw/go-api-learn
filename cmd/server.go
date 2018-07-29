package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/jcsw/go-api-learn/pkg/application"
	"github.com/jcsw/go-api-learn/pkg/infra/database"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
)

type key int

const (
	requestIDKey key = 0
)

var (
	listenAddr string
	healthy    int32
)

func main() {
	startDate := time.Now()

	flag.StringVar(&listenAddr, "listen-addr", ":8080", "server listen address")
	flag.Parse()

	logger.Info("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/", index())
	router.Handle("/health", health())
	router.HandleFunc("/customer", application.CustomerHandle)

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	mongoSession := database.CreateMongoDBSession()
	defer mongoSession.Close()

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      tracing(nextRequestID)(logging()(mongodb(mongoSession)(router))),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Info("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Info("Server is ready to handle requests at", listenAddr, "elapsed time to start was", time.Since(startDate))
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Info("Server stopped")
}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, World!")
	})
}

func health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}
			logger.Info(requestID, r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func mongodb(mongoSession *mgo.Session) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mongoSessionCopy := mongoSession.Copy()
			defer mongoSessionCopy.Close()

			ctx := context.WithValue(r.Context(), database.SessionContextKey, mongoSessionCopy)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
