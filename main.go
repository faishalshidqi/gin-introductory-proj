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
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"strconv"
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
	mongoDb := os.Getenv("MONGO_DB")
	redisUri := os.Getenv("REDIS_URI")
	redisPass := os.Getenv("REDIS_PASS")
	redisDb := os.Getenv("REDIS_DB")
	convRedisDb, _ := strconv.Atoi(redisDb)

	config = utils.ApiConfig{
		MongoURI:  mongoUri,
		MongoDB:   mongoDb,
		RedisUri:  redisUri,
		RedisPass: redisPass,
		RedisDB:   convRedisDb,
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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisUri,
		Password: config.RedisPass,
		DB:       config.RedisDB,
	})
	status := redisClient.Ping()
	log.Println(status)
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()
	router.GET("/recipes/search", recipesHandler.SearchRecipeHandler)
	router.POST("/recipes", recipesHandler.PostRecipeHandler)
	router.GET("/recipes", recipesHandler.RetrieveRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/:id", recipesHandler.RetrieveRecipeByIdHandler)
	router.Run(":9000")
}
