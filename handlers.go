package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/konfortes/go-server-utils/server"
)

// Person ...
type Person struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age"`
}

func handlers() []server.Handler {
	return []server.Handler{
		{
			// http POST localhost:3000/person name='ronen' age:=36
			Method:  http.MethodPost,
			Pattern: "/person",
			H:       personHandler,
		},
	}
}

func personHandler(c *gin.Context) {
	var person Person

	if err := c.BindJSON(&person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// functionToTrace(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{"age": person.Age})
}

// func functionToTrace(ctx context.Context) {
// 	span, _ := opentracing.StartSpanFromContext(ctx, "functionToTrace")
// 	defer span.Finish()

// 	// do some stuff.
// }
