package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	serverutils "github.com/konfortes/go-server-utils/serverutils"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	tracer *opentracing.Tracer
)

const (
	serviceName = "my-service-name"
)

func main() {
	initialize()
	router := gin.Default()

	serverutils.SetMiddlewares(router, tracer, serviceName)
	serverutils.SetRoutes(router, serviceName)
	setRoutes(router)

	srv := &http.Server{
		Addr:    ":" + serverutils.GetEnvOr("PORT", "3000"),
		Handler: router,
	}

	go func() {
		log.Println("listening on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	serverutils.GracefulShutdown(srv)
}

func initialize() {
	if serverutils.GetEnvOr("TRACING_ENABLED", "false") == "true" {
		tracer = serverutils.InitJaeger(serviceName)
	}
}
