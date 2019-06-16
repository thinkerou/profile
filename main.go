package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkerou/profile/server"
)

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	// r.LoadHTMLGlob("resources/template/*")

	r.GET("/api/user/:user", server.GetUserProfile)

	r.Run()
}

func render(c *gin.Context, data gin.H, templateName string) {
	c.JSON(http.StatusOK, data["payload"])
}
