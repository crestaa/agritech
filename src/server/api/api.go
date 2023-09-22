package api

import (
	"github.com/gin-gonic/gin"
)

func Serve() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	r.GET("/:id/test", func(c *gin.Context) {
		c.String(200, c.Param("id"))
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
