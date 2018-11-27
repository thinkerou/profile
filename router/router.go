package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thinkerou/your-profile-about-github/server"
)

func Load() http.Handler {
	r := gin.Default()

	r.GET("/api/user/:user", server.GetUserProfile)
	r.GET("/user/:user", server.RenderUserProfile)
	r.GET("/search", server.Search)
	// websocket 

	return r
}

