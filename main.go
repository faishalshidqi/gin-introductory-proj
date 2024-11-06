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
	"github.com/faishalshidqi/gin-introductory-proj/src/authentication"
	"github.com/faishalshidqi/gin-introductory-proj/src/handlers"
	"github.com/faishalshidqi/gin-introductory-proj/src/utils"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
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
var authHandler *authentication.AuthHandler

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
	collectionUsers := client.Database(config.MongoDB).Collection("users")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisUri,
		Password: config.RedisPass,
		DB:       config.RedisDB,
	})
	status := redisClient.Ping()
	log.Println(status)
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	authHandler = authentication.NewAuthHandler(ctx, collectionUsers)

}

func main() {
	router := gin.Default()
	store, _ := redisStore.NewStore(10, "tcp", config.RedisUri, config.RedisPass, []byte("secret"))
	authorized := router.Group("/")
	router.Use(sessions.Sessions("recipes_api", store))

	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.POST("/signup", authHandler.SignUpHandler)
	router.POST("/signout", authHandler.SignOutHandler)
	router.GET("/recipes", recipesHandler.RetrieveRecipesHandler)
	authorized.Use(authHandler.AuthMiddleware())
	authorized.POST("/recipes", recipesHandler.PostRecipeHandler)
	authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	authorized.GET("/recipes/search", recipesHandler.SearchRecipeHandler)
	authorized.GET("/recipes/:id", recipesHandler.RetrieveRecipeByIdHandler)
	router.Run(":9000")
}
