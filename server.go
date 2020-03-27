package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
)

// Person ...
type Person struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age"`
}

var (
	tracer            opentracing.Tracer
	tracerCloser      io.Closer
	shutdownHooks     []func()
	customMiddlewares []gin.HandlerFunc
)

func main() {
	initialize()

	router := gin.Default()

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

func initialize() {
	if isTracingEnabled() {
		initJaeger("service-name")
		customMiddlewares = append(customMiddlewares, jaegerMiddleware)
	}
}

func setMiddlewares(router *gin.Engine) {
	for _, middleware := range customMiddlewares {
		router.Use(middleware)
	}
}
