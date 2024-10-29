package main

import (
	"github.com/faishalshidqi/gin-introductory-proj/src/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	go router.POST("/recipes", handlers.PostRecipeHandler)
	go router.GET("/recipes", handlers.RetrieveRecipesHandler)
	go router.PUT("/recipes/:id", handlers.UpdateRecipeHandler)
	go router.DELETE("/recipes/:id", handlers.DeleteRecipeHandler)
	router.Run(":9000")
}
