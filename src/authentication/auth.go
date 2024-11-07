package authentication

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/auth0-community/go-auth0"
	"github.com/dgrijalva/jwt-go"
	"github.com/faishalshidqi/gin-introductory-proj/src/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"exp"`
}

// SignInHandler swagger:operation POST /signin auth signIn
//
//	Login with username and password
//	---
//	produces:
//	- application/json
//
// responses:
//
//	'200':
//		 description: Successful operation
//	'401':
//		 description: Invalid credentials
func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h := sha256.New()
	cur := handler.collection.FindOne(
		handler.ctx,
		bson.M{"username": user.Username, "password": string(h.Sum([]byte(user.Password)))},
	)
	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "user signed in"})
}

func (handler *AuthHandler) SignUpHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h := sha256.New()
	_, err := handler.collection.InsertOne(
		handler.ctx,
		bson.M{
			"username": user.Username,
			"password": string(h.Sum([]byte(user.Password))),
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user successfully created",
	})
}

// swagger:operation POST /refresh auth refresh
// Get new token in exchange for an old one
// ---
// responses:
//
//	'200':
//		 description: Successful operation
//	'400':
//		 description: Token is new and doesn't need a refresh
//	'401':
//		 description: Invalid credentials
func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if tkn == nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token isn't expired yet"})
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutput)
}

func (handler *AuthHandler) SignOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "signed out",
	})
}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth0Domain := "https://" + os.Getenv("AUTH0_DOMAIN")
		client := auth0.NewJWKClient(
			auth0.JWKClientOptions{
				URI: fmt.Sprintf("%v/.well-known/jwks.json", auth0Domain),
			},
			nil,
		)
		configuration := auth0.NewConfiguration(
			client,
			[]string{os.Getenv("AUTH0_API_IDENTIFIER")},
			fmt.Sprintf("%v/", auth0Domain),
			jose.RS256,
		)
		validator := auth0.NewValidator(configuration, nil)
		_, err := validator.ValidateRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
