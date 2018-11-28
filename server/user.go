package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thinkerou/profile/model"
)

func GetUserProfile(c *gin.Context) {
	user := c.Param("user")

}

func RenderUserProfile(c *gin.Context) {
}
