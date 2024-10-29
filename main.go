/*
Recipes API

This is a sample recipes API. You can find out more about the API at https://github.com/faishalshidqi/gin-introductory-proj
Schemes: http
Host: localhost:9000
Version: 1.0.0
Contact: Faishal Shidqi <faishalshidqi.work@proton.me>

Consumes:
- application/json

Produces:
- application/json

swagger:meta
*/
package main

import (
	"github.com/faishalshidqi/gin-introductory-proj/src/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/recipes/search", handlers.SearchRecipeHandler)
	router.POST("/recipes", handlers.PostRecipeHandler)
	router.GET("/recipes", handlers.RetrieveRecipesHandler)
	router.PUT("/recipes/:id", handlers.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", handlers.DeleteRecipeHandler)
	router.Run(":9000")
}
