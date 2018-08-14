package application

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/jcsw/go-api-learn/pkg/infra/cache"

	"github.com/jcsw/go-api-learn/pkg/infra/database"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
	"github.com/jcsw/go-api-learn/pkg/infra/properties"
)

type key int

const (
	requestIDKey key = 0
)

var healthy int32

// App define the app
type App struct {
	server *http.Server
}

// Initialize initialize the all components to app
func (app *App) Initialize(env string) {

	logger.Info("f=Initialize env=%s", env)

	properties.LoadProperties(env)
	cache.InitializeLocalCache()
	database.InitializeMongoDBSession()

	router := http.NewServeMux()
	router.Handle("/health", health())
	router.HandleFunc("/monitor", MonitorHandle)
	router.HandleFunc("/customer", CustomerHandle)

	app.server = &http.Server{
		Addr:         ":" + properties.AppProperties.ServerPort,
		Handler:      tracing()(logging()(router)),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
}

// Run initializes the application
func (app *App) Run() {

	logger.Info("Server is ready to handle requests at port %s", properties.AppProperties.ServerPort)
	atomic.StoreInt32(&healthy, 1)
	if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Could not listen on", properties.AppProperties.ServerPort, err)
	}
}

// Stop stop the application
func (app *App) Stop() {
	database.CloseMongoDBSession()

	logger.Info("Server is shutting down...")
	atomic.StoreInt32(&healthy, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	app.server.SetKeepAlivesEnabled(false)
	if err := app.server.Shutdown(ctx); err != nil {
		logger.Fatal("Could not gracefully shutdown the server, err=%v", err)
	}
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
			logger.Info("requestID=%s, method=%s path=%s remoteAddr=%s elapsedTime=%v",
				requestID, r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
		})
	}
}

func tracing() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = generateRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
