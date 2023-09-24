package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Serve() {
	r := gin.Default()

	r.Static("/web", "./static")

	apiHandler(r)

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func apiHandler(r *gin.Engine) {
	apiGroup := r.Group("/api")

	apiGroup.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "API test")
	})

	apiGroup.GET("/:id/test", func(c *gin.Context) {
		c.String(http.StatusOK, c.Param("id"))
	})

	apiGroup.GET("/fields/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		ret, _ := returnCampo(id)
		c.String(http.StatusOK, ret)
	})

	apiGroup.GET("/fields", func(c *gin.Context) {
		c.String(http.StatusOK, c.Param("id"))
	})
}

func returnCampo(id int) (string, error) {
	return "", nil
}

func returnCampi() ([]string, error) {
	return nil, nil
}
