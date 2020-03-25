package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setRoutes(router *gin.Engine) {
	// http localhost:8080/health
	router.GET("/health", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte("OK"))
	})

	// http POST localhost:8080/person name='ronen' age:=36
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
