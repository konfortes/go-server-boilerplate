package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

// Person ...
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	tracer            opentracing.Tracer
	tracerCloser      io.Closer
	shutdownHooks     []func()
	customMiddlewares []gin.HandlerFunc
)

func main() {
	router := gin.Default()

	if isTracingEnabled() {
		initJaeger("my-service-name")
		customMiddlewares = append(customMiddlewares, func(c *gin.Context) {
			log.Print("custom middleware works")
		})
	}

	setMiddlewares(router)
	setRoutes(router)

	srv := &http.Server{
		Addr:    ":" + getEnvOr("PORT", "3000"),
		Handler: router,
	}

	go func() {
		log.Println("listening on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	gracefulShutdown(srv)
}

func setMiddlewares(router *gin.Engine) {
	for _, middleware := range customMiddlewares {
		router.Use(middleware)
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
}

func setRoutes(router *gin.Engine) {
	// http localhost:8080/health
	router.GET("/health", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte("OK"))
	})

	// http POST localhost:8080/person name='ronen' age:=36
	router.POST("/person", func(c *gin.Context) {
		var person Person

		if err := c.ShouldBindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"age": person.Age})
	})
}

func getEnvOr(env, ifNotFound string) string {
	foundEnv, found := os.LookupEnv(env)

	if found {
		return foundEnv
	}

	return ifNotFound
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, hook := range shutdownHooks {
		hook()
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func initJaeger(service string) {
	cfg := &jaegerConfig.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	var err error
	tracer, tracerCloser, err = cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	shutdownHooks = append(shutdownHooks, func() {
		tracerCloser.Close()
	})

	opentracing.SetGlobalTracer(tracer)
}

func isTracingEnabled() bool {
	value := getEnvOr("TRACING_ENABLED", "false")

	return value == "true"
}
