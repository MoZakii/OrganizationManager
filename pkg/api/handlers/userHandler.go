package handlers

import (
	"MoZaki-Organization-Manager/pkg/controllers"
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"MoZaki-Organization-Manager/pkg/database/mongodb/repository"
	"MoZaki-Organization-Manager/pkg/utils"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = repository.OpenCollection(repository.Client, "user")
var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking the email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This email already exists"})
			return
		}

		password := controllers.HashPassword(*user.Password)
		user.Password = &password
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, err := utils.GenerateAllTokens(*user.Email, *user.Name, user.User_id)

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured generating tokens"})
			return
		}
		user.Token = &token
		user.Refresh_Token = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User Item was not created"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "You have successfully signed up."})

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password is incorrect"})
			return
		}

		passwordIsValid, msg := controllers.VerifyPassword(*user.Password, *foundUser.Password)

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, err := utils.GenerateAllTokens(*foundUser.Email, *foundUser.Name, foundUser.User_id)
		if utils.HasError(err, msg, c, http.StatusInternalServerError) {
			return
		}

		utils.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully",
			"access_token":  token,
			"refresh_token": refreshToken,
		})
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Input struct {
			Token_ string `json:"refresh_token"`
		}
		var in Input
		var refreshToken, token string
		if err := c.BindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		refreshToken = in.Token_

		claims, err := utils.ValidateToken(refreshToken)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		name, email, uid := claims["Name"].(string), claims["Email"].(string), claims["Uid"].(string)
		token, refreshToken, err = utils.GenerateAllTokens(name, email, uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, refreshToken = utils.UpdateAllTokens(token, refreshToken, uid)
		c.JSON(http.StatusOK, gin.H{
			"Token":         token,
			"Refresh Token": refreshToken,
			"Message":       "Token refreshed successfully",
		})
	}
}

func GetUserByEmail(c *gin.Context, email *string) (foundUser models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"email": *email}).Decode(&foundUser)
	if err != nil {
		return
	}

	return
}
