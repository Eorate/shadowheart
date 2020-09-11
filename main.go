package main

import "github.com/gin-gonic/gin"

func getMetrics() map[string]int {
	var stats = make(map[string]int)
	stats["Maintainability(mins)"] = 0
	stats["Test Coverage(%)"] = 92
	stats["Code Smells"] = 0
	stats["Duplication"] = 0
	stats["Other Issues"] = 0
	return stats
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(200, getMetrics())
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run("0.0.0.0:8080")

}
