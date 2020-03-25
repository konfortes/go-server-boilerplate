package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

func isTracingEnabled() bool {
	value := getEnvOr("TRACING_ENABLED", "false")

	return value == "true"
}

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

func jaegerMiddleware(c *gin.Context) {
	carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
	upstreamCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier)

	if err == nil {
		span := tracer.StartSpan(c.Request.URL.Path, ext.RPCServerOption(upstreamCtx))
		defer span.Finish()

		ctx := context.Background()
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.Clone(ctx)
	} else {
		log.Printf("error extracting span from request: %s", err)
	}

	c.Next()
}
