package application

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/jcsw/go-api-learn/pkg/application/handlers"
	"github.com/jcsw/go-api-learn/pkg/service"

	"github.com/jcsw/go-api-learn/pkg/infra/cache"
	"github.com/jcsw/go-api-learn/pkg/infra/cache/cachestore"
	"github.com/jcsw/go-api-learn/pkg/infra/database"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
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
	server    *http.Server
	startDate time.Time
}

// Initialize initialize the all components to app
func (app *App) Initialize(env string) {
	app.startDate = time.Now()

	logger.Info("Initialize server by env [%s]", env)

	properties.LoadProperties(env)
	cache.InitializeLocalCache()
	database.InitializeMongoClient()

	router := http.NewServeMux()
	router.HandleFunc("/health", health)

	router.HandleFunc("/monitor", handlers.MonitorHandler)

	customerRepository := repository.Repository{MongoClient: database.RetrieveMongoClient()}
	customerCacheStore := cachestore.CacheStore{}

	customerAggregate := service.CustomerAggregate{Repository: &customerRepository, CacheStore: &customerCacheStore}

	customerHandler := handlers.CustomerHandler{CAggregate: &customerAggregate}

	router.HandleFunc("/customer", customerHandler.Register)

	app.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", properties.AppProperties.ServerPort),
		Handler:      tracing()(logging()(router)),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
}

// Start initializes the application
func (app *App) Start() {
	logger.Info("Server is ready to handle requests at port %d, elapsed time to start was %v", properties.AppProperties.ServerPort, time.Since(app.startDate))

	atomic.StoreInt32(&healthy, 1)
	if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Could not listen on port [%s]\n%v", properties.AppProperties.ServerPort, err)
	}
}

// Stop stop the application
func (app *App) Stop() {
	logger.Info("Server is shutting down...")

	atomic.StoreInt32(&healthy, 0)

	database.CloseMongoClient()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.server.SetKeepAlivesEnabled(false)
	if err := app.server.Shutdown(ctx); err != nil {
		logger.Fatal("Could not gracefully shutdown the server\n%v", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
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
				requestID = newRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func newRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
