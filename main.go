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
	"context"
	"github.com/faishalshidqi/gin-introductory-proj/src/handlers"
	"github.com/faishalshidqi/gin-introductory-proj/src/models"
	"github.com/faishalshidqi/gin-introductory-proj/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

var recipes []models.Recipe
var ctx context.Context
var err error
var client *mongo.Client

func init() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
		return
	}

	mongoUri := os.Getenv("MONGO_URI")
	mongoDb := os.Getenv("MONGODB")
	config := utils.ApiConfig{
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
			collection := client.Database(config.MongoDB).Collection("recipes")
			insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Inserted %v recipes", len(insertManyResult.InsertedIDs))

	*/
}
func main() {
	router := gin.Default()
	router.GET("/recipes/search", handlers.SearchRecipeHandler)
	router.POST("/recipes", handlers.PostRecipeHandler)
	router.GET("/recipes", handlers.RetrieveRecipesHandler)
	router.PUT("/recipes/:id", handlers.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", handlers.DeleteRecipeHandler)
	router.Run(":9000")
}
