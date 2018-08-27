package application

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"gopkg.in/macaron.v1"

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
	server *macaron.Macaron
}

// Initialize initialize the all components to app
func (app *App) Initialize(env string) {
	logger.Info("Server is starting... env=%s", env)

	properties.LoadProperties(env)
	cache.InitializeLocalCache()
	database.InitializeMongoDBSession()

	app.server = macaron.New()

	app.server.Map(logger.GetConfiguredLogger())

	app.server.Use(macaron.Renderer())
	app.server.Use(macaron.Recovery())
	app.server.Use(logging)

	app.server.Before(tracing)

	app.server.Get("/health", health)
	app.server.Get("/monitor", MonitorHandle)

	app.server.Route("/customer", "GET,POST", CustomerHandle)
}

// Run initializes the application
func (app *App) Run() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		app.close()
		logger.Fatal("Server stopped")
	}()

	atomic.StoreInt32(&healthy, 1)
	app.server.Run(properties.AppProperties.ServerPort)
}

// Stop stop the application
func (app *App) close() {
	logger.Info("Server is shutting down...")
	atomic.StoreInt32(&healthy, 0)
	database.CloseMongoDBSession()
}

func health(ctx *macaron.Context) {
	if atomic.LoadInt32(&healthy) == 1 {
		ctx.Resp.WriteHeader(http.StatusOK)
		return
	}
	ctx.Resp.WriteHeader(http.StatusInternalServerError)
}

func logging(ctx *macaron.Context) {
	start := time.Now()

	ctx.Next()

	requestID := ctx.Resp.Header().Get("X-Request-Id")
	if requestID == "" {
		requestID = "unknown"
	}

	logger.Info("requestID=%s, method=%s path=%s remoteAddr=%s statusCode=%d elapsedTime=%v",
		requestID, ctx.Req.Method, ctx.Req.RequestURI, ctx.RemoteAddr(), ctx.Resp.Status(), time.Since(start))
}

func tracing(w http.ResponseWriter, r *http.Request) bool {
	requestID := r.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = newRequestID()
	}
	w.Header().Set("X-Request-Id", requestID)
	return false
}

func newRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
