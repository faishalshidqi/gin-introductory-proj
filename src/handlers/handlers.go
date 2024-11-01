package handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"

	"github.com/faishalshidqi/gin-introductory-proj/src/models"
	"github.com/gin-gonic/gin"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx:        ctx,
	}
}

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
//	'500':
//		 description: Internal server error
func (handler *RecipesHandler) PostRecipeHandler(ctx *gin.Context) {
	var recipe models.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//id := uuid.New()
	id := primitive.NewObjectID()
	pubshAt := time.Now()
	recipe.ID = id
	recipe.PublishedAt = pubshAt
	recipe.UpdatedAt = pubshAt
	_, err := handler.collection.InsertOne(ctx, recipe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed inserting a new recipe",
		})
		return
	}
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
// '500':
// description: Internal server error
func (handler *RecipesHandler) RetrieveRecipesHandler(ctx *gin.Context) {
	recipes := make([]models.Recipe, 0)
	cur, err := handler.collection.Find(ctx, bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cur, ctx)
	for cur.Next(ctx) {
		var recipe models.Recipe
		if err := cur.Decode(&recipe); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		recipes = append(recipes, recipe)
	}
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
//		 description: id not found
//	 '500':
//		 description: internal server error
func (handler *RecipesHandler) UpdateRecipeHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	var recipe models.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid id",
		})
		return
	}
	find := handler.collection.FindOne(ctx, bson.M{"_id": objectID})
	err = find.Decode(&recipe)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	_, err = handler.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.D{
			{
				"$set",
				bson.D{
					{"name", recipe.Name},
					{"instructions", recipe.Instructions},
					{"ingredients", recipe.Ingredients},
					{"tags", recipe.Tags},
					{"updatedat", time.Now()},
				},
			},
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Recipe has been updated",
		},
	)
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
//	 '400':
//		 description: Invalid id
//	 '404':
//		 description: id not found
//	 '500':
//		 description: internal server error
func (handler *RecipesHandler) DeleteRecipeHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	var recipe models.Recipe
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}
	find := handler.collection.FindOne(ctx, bson.M{"_id": objectId})
	err = find.Decode(&recipe)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	_, err = handler.collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
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
//	'500':
//		 description: Internal server error
func (handler *RecipesHandler) SearchRecipeHandler(ctx *gin.Context) {
	tag := ctx.Query("tag")
	var recipes []models.Recipe
	find, err := handler.collection.Find(ctx, bson.M{"tags": tag})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(find, ctx)
	for find.Next(ctx) {
		recipe := models.Recipe{}
		if err := find.Decode(&recipe); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		recipes = append(recipes, recipe)
	}
	ctx.JSON(http.StatusOK, recipes)
}
