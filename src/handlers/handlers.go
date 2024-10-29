package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/faishalshidqi/gin-introductory-proj/src/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PostRecipeHandler(ctx *gin.Context) {
	recipes := make([]models.Recipe, 0)
	var recipe models.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id := uuid.New()
	pubshAt := time.Now()
	recipe.ID = id
	recipe.PublishedAt = pubshAt
	recipe.UpdatedAt = pubshAt
	recipes = append(recipes, recipe)
	ctx.JSON(http.StatusOK, recipe)
}

func IndexHandler(ctx *gin.Context) {
	name := ctx.Params.ByName("name")

	ctx.JSON(200, gin.H{
		"message": fmt.Sprintf("hello %v", name),
	})
}

type Person struct {
	XMLName   xml.Name `xml:"person"`
	FirstName string   `xml:"firstName,attr"`
	LastName  string   `xml:"lastName,attr"`
}

func PersonHandler(ctx *gin.Context) {
	ctx.XML(200, Person{
		FirstName: "Tester",
		LastName:  "Testing",
	})
}
