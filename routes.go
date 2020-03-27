package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func setRoutes(router *gin.Engine) {
	// http localhost:3000/health
	router.GET("/health", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte("OK"))
	})

	// http localhost:3000/metrics
	p := ginprometheus.NewPrometheus("service_name")
	p.Use(router)

	// http POST localhost:3000/person name='ronen' age:=36
	router.POST("/person", func(c *gin.Context) {
		var person Person

		if err := c.BindJSON(&person); err != nil {
			// TODO: log to span
			return
		}

		// functionToTrace(c.Request.Context())

		c.JSON(http.StatusOK, gin.H{"age": person.Age})
	})
}

// func functionToTrace(ctx context.Context) {
// 	span, _ := opentracing.StartSpanFromContext(ctx, "functionToTrace")
// 	defer span.Finish()

// 	// do some stuff.
// }
