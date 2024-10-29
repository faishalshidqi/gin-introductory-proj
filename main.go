package main

import (
	"time"

	"github.com/faishalshidqi/gin-introductory-proj/src/handlers"
	"github.com/gin-gonic/gin"
)

type Recipe struct {
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"published_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

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
