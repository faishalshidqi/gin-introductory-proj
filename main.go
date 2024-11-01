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
	"github.com/faishalshidqi/gin-introductory-proj/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

var ctx context.Context
var err error
var client *mongo.Client
var config utils.ApiConfig
var recipesHandler *handlers.RecipesHandler

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
	ctx = context.Background()
	client, err = mongo.Connect(
		ctx,
		options.Client().ApplyURI(config.MongoURI),
	)
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(config.MongoDB).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection)
}

func main() {
	router := gin.Default()
	router.GET("/recipes/search", recipesHandler.SearchRecipeHandler)
	router.POST("/recipes", recipesHandler.PostRecipeHandler)
	router.GET("/recipes", recipesHandler.RetrieveRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.Run(":9000")
}
