package main

import (
	"bookcabin-app-go/src/libs"
	"bookcabin-app-go/src/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	libs.LoadEnv()

	router := gin.Default()
	routes.RegisterRoutes(router)

	port := libs.GetEnv("APP_PORT", "8080")
	router.Run(":" + port)
}
