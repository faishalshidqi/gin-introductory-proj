package main

import (
	"github.com/faishalshidqi/gin-introductory-proj/src/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET(
		"/:name", handlers.IndexHandler,
	)
	router.GET(
		"/person", handlers.PersonHandler,
	)
	router.Run(":9000")
}
