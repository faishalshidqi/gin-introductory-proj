package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/faishalshidqi/gin-introductory-proj/src/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []models.Recipe

func init() {
	recipes = make([]models.Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

func PostRecipeHandler(ctx *gin.Context) {
	var recipe models.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//id := uuid.New()
	id := xid.New().String()
	pubshAt := time.Now()
	recipe.ID = id
	recipe.PublishedAt = pubshAt
	recipe.UpdatedAt = pubshAt
	recipes = append(recipes, recipe)
	ctx.JSON(http.StatusOK, recipe)
}

func RetrieveRecipesHandler(ctx *gin.Context) {
	log.Println(recipes)
	ctx.JSON(http.StatusOK, recipes)
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
