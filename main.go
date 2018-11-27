package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/welcome", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "Hello world!"})
	})
	r.Run()
}
