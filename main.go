package main

import (
	"github.com/faishalshidqi/gin-introductory-proj/src/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/recipes", handlers.PostRecipeHandler)
	router.GET("/recipes", handlers.RetrieveRecipesHandler)
	router.Run(":9000")
}
