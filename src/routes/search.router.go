package routes

import (
	"bookcabin-app-go/src/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterSearchRoutes(router *gin.Engine) {
	routeGroup := router.Group("/search")
	routeGroup.POST("/", handlers.SearchFlights)
}
