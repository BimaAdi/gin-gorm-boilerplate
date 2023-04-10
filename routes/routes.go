package routes

import (
	"github.com/gin-gonic/gin"
)

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func GetRoutes(router *gin.Engine) *gin.Engine {
	noPrefixRoutes := router.Group("/")
	authRoutes(noPrefixRoutes)
	userRoutes(noPrefixRoutes)
	return router
}
