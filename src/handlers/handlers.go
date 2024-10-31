package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/faishalshidqi/gin-introductory-proj/src/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []models.Recipe

// swagger:operation POST /recipes/ recipes addRecipe
// Create a new recipe
// ---
// parameters:
//   - name: name
//     in: body
//     description: name of the recipe
//     required: true
//     type: string
//   - name: tags
//     in: body
//     description: tags of the recipe
//     required: true
//     schema:
//     type: array
//     items:
//     type: string
//   - name: ingredients
//     in: body
//     description: ingredients of the recipe
//     required: true
//     schema:
//     type: array
//     items:
//     type: string
//   - name: instructions
//     in: body
//     description: instructions to make the recipe
//     required: true
//     schema:
//     type: array
//     items:
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//		 description: Successful operation
//	'400':
//		 description: Invalid input
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

// swagger:operation GET /recipes recipes retrieveRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:

// '200':
// description: Successful operation
func RetrieveRecipesHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, recipes)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     schema:
//     type: array
//     items:
//     type: string
//   - name: name
//     in: body
//     description: name of the recipe
//     required: true
//     schema:
//     type: array
//     items:
//     type: string
//   - name: tags
//     in: body
//     description: tags of the recipe
//     required: true
//     schema:
//     type: array
//     items:
//     type: string
//   - name: ingredients
//     in: body
//     description: ingredients of the recipe
//     required: true
//     schema:
//     type: array
//     items:
//     type: string
//   - name: instructions
//     in: body
//     description: instructions to make the recipe
//     required: true
//     schema:
//     type: array
//     items:
//     type: string
//
// produces:
// - application/json
// responses:
//
//	 '200':
//	 	 description: Successful operation
//	 '400':
//	 	 description: Invalid input
//	 '404':
//		 description: Invalid recipe ID
func UpdateRecipeHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	var recipe models.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	recipe.ID = id
	recipes[index] = recipe
	ctx.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/{id} recipes removeRecipe
// Delete a recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	 '200':
//	 	 description: Successful operation
//	 '404':
//		 description: Invalid recipe ID
func DeleteRecipeHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	recipes = append(recipes[:index], recipes[index+1:]...)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

// swagger:operation GET /recipes/search recipes searchRecipe
// Look for recipe(s) with given tag
// ---
// parameters:
//   - name: tag
//     in: query
//     description: tag of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//		 description: Successful operation
func SearchRecipeHandler(ctx *gin.Context) {
	tag := ctx.Query("tag")
	listOfRecipes := make([]models.Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}
	ctx.JSON(http.StatusOK, listOfRecipes)
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
