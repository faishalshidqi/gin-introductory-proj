package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/faishalshidqi/gin-introductory-proj/src/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/faishalshidqi/gin-introductory-proj/src/models"
	"github.com/gin-gonic/gin"
)

var recipes []models.Recipe
var ctx context.Context
var err error
var client *mongo.Client
var config utils.ApiConfig

func init() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
		return
	}

	mongoUri := os.Getenv("MONGO_URI")
	mongoDb := os.Getenv("MONGODB")
	config = utils.ApiConfig{
		MongoURI: mongoUri,
		MongoDB:  mongoDb,
	}
	/*
		just in case reloading recipes.json is needed
		recipes = make([]models.Recipe, 0)
		file, _ := os.ReadFile("recipes.json")
		_ = json.Unmarshal(file, &recipes)
	*/
	ctx = context.Background()
	client, err = mongo.Connect(
		ctx,
		options.Client().ApplyURI(config.MongoURI),
	)
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	/*
		just in case reloading recipes.json is needed
			var listOfRecipes []interface{}
			for _, recipe := range recipes {
				listOfRecipes = append(listOfRecipes, recipe)
			}
			insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Inserted %v recipes", len(insertManyResult.InsertedIDs))

	*/
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
func PostRecipeHandler(ctx *gin.Context) {
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
	collection := client.Database(config.MongoDB).Collection("recipes")
	_, err := collection.InsertOne(ctx, recipe)
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
func RetrieveRecipesHandler(ctx *gin.Context) {
	collection := client.Database(config.MongoDB).Collection("recipes")

	cur, err := collection.Find(ctx, bson.M{})
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
	collection := client.Database(config.MongoDB).Collection("recipes")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid id",
		})
		return
	}
	_, err = collection.UpdateOne(
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
//	 '404':
//		 description: Invalid recipe ID
func DeleteRecipeHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}
	collection := client.Database(config.MongoDB).Collection("recipes")
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectId})
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
func SearchRecipeHandler(ctx *gin.Context) {
	tag := ctx.Query("tag")
	collection := client.Database(config.MongoDB).Collection("recipes")
	var recipes []models.Recipe
	find, err := collection.Find(ctx, bson.M{"tags": tag})
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
