package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(200, BuildRepositoryStats())
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run("0.0.0.0:8080")

}
